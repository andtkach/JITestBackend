package main

import (
	"fmt"
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// Define our Handler
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "Invalid Request", fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("Successfully called by - %s", event.Username), nil
}

func main() {
	lambdaApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return lambdaApp.ApiHandler.RegisterUser(request)
		case "/login":
			return lambdaApp.ApiHandler.LoginUser(request)
		case "/me":
			return middleware.ValidateJWTMiddleware(lambdaApp.ApiHandler.GetUser)(request)
		case "/role":
			return middleware.ValidateJWTMiddleware(lambdaApp.ApiHandler.UpdateRole)(request)
		case "/list":
			return middleware.ValidateJWTMiddleware(lambdaApp.ApiHandler.ListUsers)(request)
		case "/remove":
			return middleware.ValidateJWTMiddleware(lambdaApp.ApiHandler.RemoveUser)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}

	})
}
