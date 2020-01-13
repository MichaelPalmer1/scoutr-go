package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/MichaelPalmer1/simple-api-go/config"
	"github.com/MichaelPalmer1/simple-api-go/filterbuilder"
	"github.com/MichaelPalmer1/simple-api-go/models"
	"github.com/MichaelPalmer1/simple-api-go/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Record : Item in Dynamo
type Record map[string]interface{}

// Initialize - Creates connection to DynamoDB
func Initialize(config *config.Config) *dynamodb.DynamoDB {
	usr, _ := user.Current()
	creds := credentials.NewSharedCredentials(filepath.Join(usr.HomeDir, ".aws/credentials"), "default")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	}))

	svc := dynamodb.New(sess)

	return svc
}

func main() {
	// Command line arguments
	var config config.Config
	flag.StringVar(&config.DataTable, "data-table", "", "Data table")
	flag.StringVar(&config.AuthTable, "auth-table", "", "Auth table")
	flag.StringVar(&config.GroupTable, "group-table", "", "Group table")
	flag.StringVar(&config.AuditTable, "audit-table", "", "Audit table")
	flag.IntVar(&config.LogRetentionDays, "log-retention-days", 30, "Days to retain read logs")
	flag.Parse()

	svc := Initialize(&config)

	user, err := utils.GetUser("michael", config.AuthTable, config.GroupTable, svc, nil, nil)
	if err != nil {
		panic(err)
	}

	input := &dynamodb.ScanInput{
		TableName: aws.String(config.DataTable),
	}

	output, err := svc.Scan(input)
	if err != nil {
		fmt.Println("encountered error", err)
		return
	}

	records := []models.Record{}
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &records)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal items, %v", err))
	}

	out, err := json.Marshal(records)
	fmt.Println(string(out))
	records = utils.PostProcess(records, user)
	out, err = json.Marshal(records)
	fmt.Println(string(out))

	filters := make(map[string]string)
	filters["item__contains"] = "test"

	filterbuilder.BuildFilter(user, filters)

}
