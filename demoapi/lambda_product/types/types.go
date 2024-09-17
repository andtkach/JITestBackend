package types

import (
	"lambda-func/common"
)

type Product struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Manager     string `json:"manager"`
}

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type UpdateProductRequest struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type UserContext struct {
	Username string
	Role     string
}

type ProductResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

func NewProduct(productRequest CreateProductRequest, manager string) (Product, error) {
	return Product{
		Id:          common.GenerateStrignID(),
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		Manager:     manager,
	}, nil
}
