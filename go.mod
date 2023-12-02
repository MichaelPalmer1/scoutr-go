module github.com/MichaelPalmer1/scoutr-go

go 1.20

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2 v1.23.5
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.12.9
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.6.9
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.35.2
	github.com/aws/aws-sdk-go-v2/service/cloudtraildata v1.5.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.26.3
	github.com/aws/smithy-go v1.18.1
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.8.9 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
