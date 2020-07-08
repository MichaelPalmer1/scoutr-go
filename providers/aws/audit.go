package aws

import (
	"sort"
	"time"

	"github.com/MichaelPalmer1/scoutr-go/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	log "github.com/sirupsen/logrus"
)

// ListAuditLogs : List audit logs
func (api DynamoAPI) ListAuditLogs(req models.Request, pathParams map[string]string, queryParams map[string]string) ([]models.AuditLog, error) {
	// Only fetch audit logs if the table is configured
	if api.Config.AuditTable == "" {
		return nil, &models.NotFound{
			Message: "Audit logs are not enabled",
		}
	}

	// Get the user
	_, err := api.InitializeRequest(api, req)
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
		queryParams[key] = value
	}

	// Build filters
	rawConds, hasConditions, err := api.Filter(&api.Filtering, nil, queryParams)
	if err != nil {
		log.Errorln("Error encountered during filtering", err)
		return nil, err
	}
	if hasConditions {
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

	// Download the data
	data, err := scanAudit(&input, api.Client)
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
func (api DynamoAPI) auditLog(action string, request models.Request, user models.User, resource *map[string]string, changes *map[string]string) {
	// Only send audit logs if the table is configured
	if api.Config.AuditTable == "" {
		return
	}

	// Create audit log
	now := time.Now().UTC()
	auditLog := models.AuditLog{
		Time: now.Format(time.RFC3339Nano),
		User: models.AuditUser{
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

	// Add expiry time for read events
	if action == "GET" || action == "LIST" || action == "SEARCH" {
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
		auditLog.Body = *changes
	}

	// Add resource
	if resource != nil {
		auditLog.Resource = *resource
	}

	// Marshal the audit log to Dynamo format
	item, err := dynamodbattribute.MarshalMap(auditLog)
	if err != nil {
		log.Errorln("Failed to marshal the audit log", err)
		log.Infof("Failed audit log: '%v'", auditLog)
		return
	}

	// Generate the put item input
	input := dynamodb.PutItemInput{
		TableName: aws.String(api.Config.AuditTable),
		Item:      item,
	}

	// Add the record to dynamo
	_, err = api.Client.PutItem(&input)
	if err != nil {
		log.Errorln("Failed to put audit log in Dynamo", err)
		log.Infof("Failed audit log: '%v'", auditLog)
		return
	}
}
