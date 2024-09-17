package main

import (
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambdaApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/list":
			return lambdaApp.ApiHandler.ListProducts(request)
		case "/one":
			return lambdaApp.ApiHandler.GetProduct(request)
		case "/create":
			return middleware.ValidateJWTMiddleware(lambdaApp.ApiHandler.CreateProduct)(request)
		case "/update":
			return middleware.ValidateJWTMiddleware(lambdaApp.ApiHandler.UpdateProduct)(request)
		case "/delete":
			return middleware.ValidateJWTMiddleware(lambdaApp.ApiHandler.DeleteProduct)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}

	})
}
