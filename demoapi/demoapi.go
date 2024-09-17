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

	tableProducts := awsdynamodb.NewTable(stack, jsii.String(common.ProductTableName), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName:     jsii.String(common.ProductTableName),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	functionUsers := awslambda.NewFunction(stack, jsii.String(common.UserFunctionName), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda_user/user_function.zip"), nil),
		Handler: jsii.String("main"),
	})

	functionProducts := awslambda.NewFunction(stack, jsii.String(common.ProductFunctionName), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda_product/product_function.zip"), nil),
		Handler: jsii.String("main"),
	})

	tableUsers.GrantReadWriteData(functionUsers)
	queue.GrantSendMessages(functionUsers)

	tableProducts.GrantReadWriteData(functionProducts)

	apiUser := awsapigateway.NewRestApi(stack, jsii.String(common.UserGatewayName), &awsapigateway.RestApiProps{
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

	integrationUser := awsapigateway.NewLambdaIntegration(functionUsers, nil)

	registerResource := apiUser.Root().AddResource(jsii.String("register"), nil)
	registerResource.AddMethod(jsii.String("POST"), integrationUser, nil)

	loginResource := apiUser.Root().AddResource(jsii.String("login"), nil)
	loginResource.AddMethod(jsii.String("POST"), integrationUser, nil)

	meResource := apiUser.Root().AddResource(jsii.String("me"), nil)
	meResource.AddMethod(jsii.String("GET"), integrationUser, nil)

	roleResource := apiUser.Root().AddResource(jsii.String("role"), nil)
	roleResource.AddMethod(jsii.String("PUT"), integrationUser, nil)

	removeResource := apiUser.Root().AddResource(jsii.String("remove"), nil)
	removeResource.AddMethod(jsii.String("DELETE"), integrationUser, nil)

	listResource := apiUser.Root().AddResource(jsii.String("list"), nil)
	listResource.AddMethod(jsii.String("GET"), integrationUser, nil)

	apiProduct := awsapigateway.NewRestApi(stack, jsii.String(common.ProductGatewayName), &awsapigateway.RestApiProps{
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

	integrationProduct := awsapigateway.NewLambdaIntegration(functionProducts, nil)

	productListResource := apiProduct.Root().AddResource(jsii.String("list"), nil)
	productListResource.AddMethod(jsii.String("GET"), integrationProduct, nil)

	productOneResource := apiProduct.Root().AddResource(jsii.String("one"), nil)
	productOneResource.AddMethod(jsii.String("GET"), integrationProduct, nil)

	productCreateResource := apiProduct.Root().AddResource(jsii.String("create"), nil)
	productCreateResource.AddMethod(jsii.String("POST"), integrationProduct, nil)

	productUpdateResource := apiProduct.Root().AddResource(jsii.String("update"), nil)
	productUpdateResource.AddMethod(jsii.String("PUT"), integrationProduct, nil)

	productDeleteResource := apiProduct.Root().AddResource(jsii.String("delete"), nil)
	productDeleteResource.AddMethod(jsii.String("DELETE"), integrationProduct, nil)

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

func env() *awscdk.Environment {
	return nil
}
