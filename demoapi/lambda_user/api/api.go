package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/common"
	"lambda-func/database"
	"lambda-func/queue"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore  database.UserStore
	msgQueue queue.MessageQueue
}

func NewApiHandler(dbStore database.UserStore, msgQueue queue.MessageQueue) ApiHandler {
	return ApiHandler{
		dbStore:  dbStore,
		msgQueue: msgQueue,
	}
}

func (api ApiHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("register user request fields cannot be empty")
	}

	doesUserExist, err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("checking if user exists error %w", err)
	}

	if doesUserExist {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusConflict,
		}, nil
	}

	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal Server error",
			StatusCode: http.StatusConflict,
		}, fmt.Errorf("error hashing user password %w", err)
	}

	err = api.dbStore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user into the database %w", err)
	}

	queueMessageBody := "New user created " + user.Username
	err = api.msgQueue.SendMessage(queueMessageBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user into the database %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "Success",
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid login credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}
	accessToken := types.CreateToken(user)
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) GetUser(request events.APIGatewayProxyRequest, userContext types.UserContext) (events.APIGatewayProxyResponse, error) {

	if userContext.Username == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Anauthorized error",
			StatusCode: http.StatusUnauthorized,
		}, fmt.Errorf("user context is nil")
	}

	username := userContext.Username

	user, err := api.dbStore.GetUser(username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	successMsg := fmt.Sprintf(`{"username": "%s", "role": "%s"}`, username, user.Role)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) UpdateRole(request events.APIGatewayProxyRequest, userContext types.UserContext) (events.APIGatewayProxyResponse, error) {

	result, err := checkAdmin(userContext)
	if err != nil {
		return result, err
	}

	type RoleRequest struct {
		Username string `json:"username"`
		NewRole  string `json:"newrole"`
	}

	var roleRequest RoleRequest

	err = json.Unmarshal([]byte(request.Body), &roleRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(roleRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if !(roleRequest.NewRole == common.RoleAdmin || roleRequest.NewRole == common.RoleUser) {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Invalid role %s", roleRequest.NewRole),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user.Role = roleRequest.NewRole

	err = api.dbStore.UpdateUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	successMsg := fmt.Sprintf(`{"username": "%s", "role": "%s"}`, user.Username, user.Role)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) RemoveUser(request events.APIGatewayProxyRequest, userContext types.UserContext) (events.APIGatewayProxyResponse, error) {

	result, err := checkAdmin(userContext)
	if err != nil {
		return result, err
	}

	username := request.QueryStringParameters["username"]

	user, err := api.dbStore.GetUser(username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	err = api.dbStore.DeleteUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	successMsg := fmt.Sprintf(`user removed`, username)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) ListUsers(request events.APIGatewayProxyRequest, userContext types.UserContext) (events.APIGatewayProxyResponse, error) {

	result, err := checkAdmin(userContext)
	if err != nil {
		return result, err
	}

	users, err := api.dbStore.ListUsers()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	var userResponse []types.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, types.UserResponse{
			Username: user.Username,
			Role:     user.Role,
		})
	}

	jsonResponse, err := json.Marshal(userResponse)
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

	fmt.Println("in checkAdmin. userContext.Role is ", userContext.Role)

	if userContext.Username == "" {
		fmt.Println("in checkAdmin. userContext Username is nil")
		return events.APIGatewayProxyResponse{
			Body:       "Anauthorized error",
			StatusCode: http.StatusUnauthorized,
		}, fmt.Errorf("user context is nil")
	}

	if userContext.Role != common.RoleAdmin {
		fmt.Println("in checkAdmin. userContext Role is not admin, it is ", userContext.Role)
		return events.APIGatewayProxyResponse{
			Body:       "Anauthorized error",
			StatusCode: http.StatusUnauthorized,
		}, fmt.Errorf("user does not have enogth privileges")
	}

	return events.APIGatewayProxyResponse{}, nil
}
