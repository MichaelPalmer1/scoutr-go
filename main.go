package main

import (
	"flag"
	"fmt"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Record : Item in Dynamo
type Record struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Value string `json:"value"` 
}

func main() {
	// Command line arguments
	var config config.Config
	flag.StringVar(&config.DataTable, "data-table", "", "Data table")
	flag.StringVar(&config.AuthTable, "auth-table", "", "Auth table")
	flag.StringVar(&config.AuditTable, "audit-table", "", "Audit table")
	flag.IntVar(&config.LogRetentionDays, "log-retention-days", 30, "Days to retain read logs")
	flag.Parse()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("~/.aws/credentials", "default"),
	})

	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String(config.DataTable),
	}

	output, err := svc.Scan(input)
	if err != nil {
		fmt.Println("encountered error", err)
		return
	}

	records := []Record{}

	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &records)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal items, %v", err))
	}

	fmt.Println(records[0].Name)

}
