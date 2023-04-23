module github.com/MichaelPalmer1/scoutr-go

go 1.20

require (
	github.com/aws/aws-lambda-go v1.37.0
	github.com/aws/aws-sdk-go-v2 v1.17.8
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.15
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.4.42
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.24.6
	github.com/aws/aws-sdk-go-v2/service/cloudtraildata v1.0.4
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.18.5
	github.com/aws/smithy-go v1.13.5
	github.com/cenkalti/backoff/v4 v4.2.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.14.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.23 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)
