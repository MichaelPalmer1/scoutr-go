package aws

// Patch : Patch an item
// func (api DynamoAPI) Patch(request types.Request, partitionKey map[string]interface{}, item map[string]interface{}, validation map[string]types.FieldValidation, requiredFields []string, auditAction string) (interface{}, error) {
// 	builder := expression.NewBuilder()

// 	// Get the user
// 	user, err := api.InitializeRequest(request)
// 	if err != nil {
// 		// Bad user - pass the error through
// 		return nil, err
// 	}

// 	// Run data validation
// 	if validation != nil {
// 		log.Infoln("Running field validation")
// 		err := api.ValidateFields(validation, requiredFields, item, nil)
// 		if err != nil {
// 			log.Errorln("Field validation error", err)
// 			return nil, err
// 		}
// 	}

// 	// Build update expression
// 	var updateExpr expression.UpdateBuilder
// 	hasUpdateConds := false
// 	for key, value := range item {
// 		if !hasUpdateConds {
// 			updateExpr = expression.Set(expression.Name(key), expression.Value(value))
// 			hasUpdateConds = true
// 		} else {
// 			updateExpr = updateExpr.Set(expression.Name(key), expression.Value(value))
// 		}
// 	}
// 	builder = builder.WithUpdate(updateExpr)

// 	// Build partition key
// 	dynamoKeyParts, err := attributevalue.MarshalMap(partitionKey)
// 	if err != nil {
// 		log.Errorln("Failed to marshal partition key", err)
// 		return nil, err
// 	}

// 	// Build filters
// 	conditions, err := api.filtering.Filter(user, nil, base.FilterActionUpdate)
// 	if err != nil {
// 		log.Errorln("Error encountered during filtering", err)
// 		return nil, err
// 	}

// 	// Get key schema
// 	keySchema, err := api.Client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
// 		TableName: aws.String(api.Config.DataTable),
// 	})
// 	if err != nil {
// 		log.Errorln("Failed to describe table", err)
// 		return nil, err
// 	}

// 	// Append key schema conditions
// 	for _, schema := range keySchema.Table.KeySchema {
// 		condition := expression.Name(*schema.AttributeName).AttributeExists()
// 		conditions = api.filtering.And(conditions, condition)
// 	}

// 	// Add conditions
// 	if conditions != nil {
// 		builder = builder.WithCondition(conditions.(expression.ConditionBuilder))
// 	}

// 	// Build expression
// 	expr, err := builder.WithUpdate(updateExpr).Build()
// 	if err != nil {
// 		log.Errorln("Encountered error while building expression", err)
// 		return nil, err
// 	}

// 	// Update the item in dynamo
// 	updatedItem, err := UpdateItem[types.Record](api.Client, api.Config.DataTable, dynamoKeyParts, expr)
// 	if err != nil {
// 		log.Errorln("Error while attempting to update item in dynamo", err)

// 		// Check if this was a conditional check failure
// 		var apiErr smithy.APIError
// 		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ConditionalCheckFailedException" {
// 			return nil, &types.BadRequest{
// 				Message: "Item does not exist or you do not have permission to update it",
// 			}
// 		}

// 		return nil, err
// 	}

// 	// Create audit log
// 	api.auditLog(auditAction, request, user, partitionKey, item)

// 	return updatedItem, nil
// }
