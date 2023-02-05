package aws

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/MichaelPalmer1/scoutr-go/pkg/providers/base"
	"github.com/MichaelPalmer1/scoutr-go/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/cloudtraildata"
	cloudTrailDataTypes "github.com/aws/aws-sdk-go-v2/service/cloudtraildata/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// ListAuditLogs : List audit logs
func (api DynamoAPI) ListAuditLogs(req types.Request, pathParams map[string]string, queryParams map[string][]string) ([]types.AuditLog, error) {
	// Only fetch audit logs if the table is configured
	if api.Config.AuditTable == "" {
		return nil, &types.NotFound{
			Message: "Audit logs are not enabled",
		}
	}

	// Get the user
	_, err := api.InitializeRequest(req)
	if err != nil {
		// Bad user - pass the error through
		return nil, err
	}

	input := dynamodb.ScanInput{
		TableName: aws.String(api.Config.AuditTable),
	}

	// Generate dynamic search
	searchKey, hasSearchKey := pathParams["search_key"]
	searchValue, hasSearchValue := pathParams["search_value"]
	if hasSearchKey && hasSearchValue {
		// Map the search key and value into path params
		pathParams[searchKey] = searchValue
		delete(pathParams, "search_key")
		delete(pathParams, "search_value")
	}

	// Merge pathParams into queryParams
	for key, value := range pathParams {
		queryParams[key] = append(queryParams[key], value)
	}

	// Build filters
	if rawConds, err := api.filtering.Filter(nil, queryParams, ""); err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	} else if rawConds != nil {
		conditions := rawConds.(expression.ConditionBuilder)
		expr, err := expression.NewBuilder().WithFilter(conditions).Build()
		if err != nil {
			return nil, err
		}

		// Update scan input
		input.FilterExpression = expr.Filter()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}

	// f, err := api.cloudTrailClient.StartQuery(context.TODO(), &cloudtrail.StartQueryInput{
	// 	QueryStatement: aws.String("SELECT *"),
	// })

	// r, err := api.cloudTrailClient.DescribeQuery(context.TODO(), &cloudtrail.DescribeQueryInput{
	// 	QueryId: f.QueryId,
	// })

	// g, err := api.cloudTrailClient.GetQueryResults(context.TODO(), &cloudtrail.GetQueryResultsInput{
	// 	QueryId: f.QueryId,
	// })

	// Download the data
	data, err := Scan[types.AuditLog](api.Client, &input)
	if err != nil {
		log.Errorln("Error while attempting to list records", err)
		return nil, err
	}

	// Sort the results
	sort.Slice(data, func(i, j int) bool {
		return data[i].Time > data[j].Time
	})

	return data, nil
}

// auditLog : Creates an audit log
func (api DynamoAPI) auditLog(action string, request types.Request, user *types.User, resource map[string]interface{}, changes map[string]interface{}) {
	// Only send audit logs if the table is configured
	if api.Config.AuditTable == "" {
		return
	}

	// Create audit log
	now := time.Now().UTC()
	auditLog := types.AuditLog{
		Time: now.Format(time.RFC3339Nano),
		User: types.AuditUser{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			SourceIP:  request.SourceIP,
			UserAgent: request.UserAgent,
		},
		Action: action,
		Method: request.Method,
		Path:   request.Path,
	}

	auditEvent := types.AuditEvent{
		EventData: types.AuditEventData{
			Version: "1.0",
			UserIdentity: types.AuditEventUserIdentity{
				Type:        "ScoutrUser",
				PrincipalId: user.ID,
				Details:     user,
			},
			UserAgent:         request.UserAgent,
			SourceIPAddress:   request.SourceIP,
			EventSource:       "scoutr",
			EventName:         action,
			EventTime:         time.Now().UTC(),
			UID:               "RANDOM",
			RequestParameters: resource,
			AdditionalEventData: map[string]interface{}{
				"request": request,
				"changes": changes,
			},
		},
	}

	// Add expiry time for read events
	if action == base.AuditActionGet || action == base.AuditActionList || action == base.AuditActionSearch {
		auditLog.ExpireTime = now.AddDate(0, 0, api.Config.LogRetentionDays).Unix()
	}

	// Add query params
	if len(request.QueryParams) > 0 {
		auditLog.QueryParams = request.QueryParams
	}

	// Add body
	if request.Body != nil {
		auditLog.Body = request.Body
	} else if changes != nil {
		auditLog.Body = changes
	}

	// Add resource
	if resource != nil {
		auditLog.Resource = resource
	}

	// Marshal the audit log to Dynamo format
	// item, err := attributevalue.MarshalMap(auditLog)
	// if err != nil {
	// 	log.Errorln("Failed to marshal the audit log", err)
	// 	log.Infof("Failed audit log: '%v'", auditLog)
	// 	return
	// }

	// Generate the put item input
	// input := dynamodb.PutItemInput{
	// 	TableName: aws.String(api.Config.AuditTable),
	// 	Item:      item,
	// }

	bs, err := json.Marshal(auditEvent)
	if err != nil {
		log.WithError(err).WithField("AuditEvent", auditEvent).Errorf("Failed to marshal audit event")
		return
	}

	auditEvents := []cloudTrailDataTypes.AuditEvent{
		{
			Id:        aws.String("unique-id"),
			EventData: aws.String(string(bs)),
		},
	}

	result, err := api.auditClient.PutAuditEvents(context.TODO(), &cloudtraildata.PutAuditEventsInput{
		ChannelArn:  aws.String("channel-arn"),
		AuditEvents: auditEvents,
	})

	for _, item := range result.Failed {
		log.Errorf("Failed to record event %s - %s: %s", aws.ToString(item.Id), aws.ToString(item.ErrorCode), aws.ToString(item.ErrorMessage))
	}

	// Add the record to dynamo
	// _, err = api.Client.PutItem(context.TODO(), &input)
	// if err != nil {
	// 	log.Errorln("Failed to put audit log in Dynamo", err)
	// 	log.Infof("Failed audit log: '%v'", auditLog)
	// 	return
	// }

}
