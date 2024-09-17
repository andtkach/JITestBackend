package main

import (
	"demoapi/common"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DemoapiStackProps struct {
	awscdk.StackProps
}

func NewDemoapiStack(scope constructs.Construct, id string, props *DemoapiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	queue := awssqs.NewQueue(stack, jsii.String(common.QueueName), &awssqs.QueueProps{
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
		QueueName:         jsii.String(common.QueueName),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	})

	tableUsers := awsdynamodb.NewTable(stack, jsii.String(common.UserTableName), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName:     jsii.String(common.UserTableName),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	function := awslambda.NewFunction(stack, jsii.String(common.FunctionName), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	tableUsers.GrantReadWriteData(function)
	queue.GrantSendMessages(function)

	api := awsapigateway.NewRestApi(stack, jsii.String(common.GatewayName), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel:     awsapigateway.MethodLoggingLevel_INFO,
			DataTraceEnabled: jsii.Bool(true),
		},
		EndpointConfiguration: &awsapigateway.EndpointConfiguration{
			Types: &[]awsapigateway.EndpointType{awsapigateway.EndpointType_REGIONAL},
		},
		CloudWatchRole: jsii.Bool(true),
	})

	integration := awsapigateway.NewLambdaIntegration(function, nil)

	registerResource := api.Root().AddResource(jsii.String("register"), nil)
	registerResource.AddMethod(jsii.String("POST"), integration, nil)

	loginResource := api.Root().AddResource(jsii.String("login"), nil)
	loginResource.AddMethod(jsii.String("POST"), integration, nil)

	meResource := api.Root().AddResource(jsii.String("me"), nil)
	meResource.AddMethod(jsii.String("GET"), integration, nil)

	roleResource := api.Root().AddResource(jsii.String("role"), nil)
	roleResource.AddMethod(jsii.String("PUT"), integration, nil)

	removeResource := api.Root().AddResource(jsii.String("remove"), nil)
	removeResource.AddMethod(jsii.String("DELETE"), integration, nil)

	listResource := api.Root().AddResource(jsii.String("list"), nil)
	listResource.AddMethod(jsii.String("GET"), integration, nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewDemoapiStack(app, common.StackName, &DemoapiStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
