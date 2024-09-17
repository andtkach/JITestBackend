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

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User, error)
	UpdateUser(user types.User) error
	DeleteUser(user types.User) error
	ListUsers() ([]types.User, error)
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

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u DynamoDBClient) InsertUser(user types.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(common.UserTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
			"role": {
				S: aws.String(user.Role),
			},
		},
	}

	_, err := u.databaseStore.PutItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) UpdateUser(user types.User) error {

	update := expression.Set(expression.Name("role"), expression.Value(user.Role))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return err
	}

	item := &dynamodb.UpdateItemInput{
		TableName: aws.String(common.UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	_, err = u.databaseStore.UpdateItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) GetUser(username string) (types.User, error) {
	var user types.User

	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(common.UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u DynamoDBClient) DeleteUser(user types.User) error {

	item := &dynamodb.DeleteItemInput{
		TableName: aws.String(common.UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
		},
	}

	_, err := u.databaseStore.DeleteItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) ListUsers() ([]types.User, error) {
	var users []types.User

	result, err := u.databaseStore.Scan(&dynamodb.ScanInput{
		TableName: aws.String(common.UserTableName),
	})

	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		item := types.User{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return nil, err
		}

		users = append(users, item)
	}

	return users, nil
}
