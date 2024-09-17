package database

import (
	"fmt"
	"lambda-func/common"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ProductStore interface {
	ListProducts() ([]types.Product, error)
	GetProduct(id string) (types.Product, error)
	CreateProduct(product types.Product) error
	UpdateProduct(product types.Product) error
	DeleteProduct(product types.Product) error
}

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDB() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return DynamoDBClient{
		databaseStore: db,
	}
}

func (p DynamoDBClient) CreateProduct(product types.Product) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(common.ProductTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(product.Id),
			},
			"name": {
				S: aws.String(product.Name),
			},
			"description": {
				S: aws.String(product.Description),
			},
			"price": {
				N: aws.String(fmt.Sprintf("%d", product.Price)),
			},
			"manager": {
				S: aws.String(product.Manager),
			},
		},
	}

	_, err := p.databaseStore.PutItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (p DynamoDBClient) UpdateProduct(product types.Product) error {

	update := expression.Set(expression.Name("name"), expression.Value(product.Name))
	update = update.Set(expression.Name("description"), expression.Value(product.Description))
	update = update.Set(expression.Name("price"), expression.Value(product.Price))
	update = update.Set(expression.Name("manager"), expression.Value(product.Manager))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return err
	}

	item := &dynamodb.UpdateItemInput{
		TableName: aws.String(common.ProductTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(product.Id),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	_, err = p.databaseStore.UpdateItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (p DynamoDBClient) DeleteProduct(product types.Product) error {

	item := &dynamodb.DeleteItemInput{
		TableName: aws.String(common.ProductTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(product.Id),
			},
		},
	}

	_, err := p.databaseStore.DeleteItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (p DynamoDBClient) GetProduct(id string) (types.Product, error) {
	var product types.Product

	result, err := p.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.ProductTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return product, err
	}

	if result.Item == nil {
		return product, fmt.Errorf("product not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &product)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (p DynamoDBClient) ListProducts() ([]types.Product, error) {
	var products []types.Product

	result, err := p.databaseStore.Scan(&dynamodb.ScanInput{
		TableName: aws.String(common.ProductTableName),
	})

	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		item := types.Product{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return nil, err
		}

		products = append(products, item)
	}

	return products, nil
}
