package app

import (
	"lambda-func/api"
	"lambda-func/database"
	"lambda-func/queue"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	db := database.NewDynamoDB()
	q := queue.NewSqsClient()
	apiHandler := api.NewApiHandler(db, q)

	return App{
		ApiHandler: apiHandler,
	}
}
