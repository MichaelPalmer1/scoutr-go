package endpoints

import (
	"time"

	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
)

// AuditLog : Creates an audit log
func (api *SimpleAPI) auditLog(action string, request models.Request, user models.User, resource *map[string]string, changes *map[string]string) error {
	// Only send audit logs if the table is configured
	if api.Config.AuditTable == "" {
		return nil
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
		return err
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
		log.Infof("Failed audit log: '%s'", auditLog)
		return err
	}

	return nil
}
