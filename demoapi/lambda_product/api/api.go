package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/common"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.ProductStore
}

func NewApiHandler(dbStore database.ProductStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) CreateProduct(request events.APIGatewayProxyRequest, userContext types.UserContext) (events.APIGatewayProxyResponse, error) {
	var createProduct types.CreateProductRequest

	result, err := checkAdmin(userContext)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(request.Body), &createProduct)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if createProduct.Name == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("create productrequest is invalid")
	}

	product, err := types.NewProduct(createProduct, userContext.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server error",
			StatusCode: http.StatusConflict,
		}, fmt.Errorf("error creating database product %w", err)
	}

	err = api.dbStore.CreateProduct(product)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting product into the database %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       product.Id,
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) GetProduct(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	productId := request.QueryStringParameters["id"]

	product, err := api.dbStore.GetProduct(productId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	jsonResponse, err := json.Marshal(product)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponse),
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) UpdateProduct(request events.APIGatewayProxyRequest, userContext types.UserContext) (events.APIGatewayProxyResponse, error) {

	result, err := checkAdmin(userContext)
	if err != nil {
		return result, err
	}

	var updateProductRequest types.UpdateProductRequest

	err = json.Unmarshal([]byte(request.Body), &updateProductRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	product, err := api.dbStore.GetProduct(updateProductRequest.Id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	product.Name = updateProductRequest.Name
	product.Description = updateProductRequest.Description
	product.Price = updateProductRequest.Price

	err = api.dbStore.UpdateProduct(product)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	jsonResponse, err := json.Marshal(product)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponse),
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) DeleteProduct(request events.APIGatewayProxyRequest, userContext types.UserContext) (events.APIGatewayProxyResponse, error) {

	result, err := checkAdmin(userContext)
	if err != nil {
		return result, err
	}

	productId := request.QueryStringParameters["id"]

	product, err := api.dbStore.GetProduct(productId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	err = api.dbStore.DeleteProduct(product)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	successMsg := fmt.Sprintf(`product removed`, productId)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) ListProducts(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	products, err := api.dbStore.ListProducts()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	var productResponse []types.ProductResponse
	for _, product := range products {
		productResponse = append(productResponse, types.ProductResponse{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	jsonResponse, err := json.Marshal(productResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonResponse),
		StatusCode: http.StatusOK,
	}, nil
}

func checkAdmin(userContext types.UserContext) (events.APIGatewayProxyResponse, error) {
	if userContext.Username == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Anauthorized error",
			StatusCode: http.StatusUnauthorized,
		}, fmt.Errorf("user context is nil")
	}

	if userContext.Role != common.RoleAdmin {
		return events.APIGatewayProxyResponse{
			Body:       "Anauthorized error",
			StatusCode: http.StatusUnauthorized,
		}, fmt.Errorf("user does not have enogth privileges")
	}

	return events.APIGatewayProxyResponse{}, nil
}
