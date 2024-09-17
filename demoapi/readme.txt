
API Gayeway -> AWS Lambda Function -> AWS DynamoDB
AWS CloudFormation

Curl as a client
Deployment with a scripts
Infrastructure as a code


DB:
- Users
- Products

API
- Register new user
- Login user
- CRUD for products

AWS on andrii0spain@gmail.com account (eu-central-1 region)

user: aws-test-user
group: aws-test-group
Access key: ACCESS-KEY
Secret access key: ACCESS-SECRET
 
Install tools:
1. Install Go
2. Install AWS CLI
aws --version
3. Create IAM user and configure user with access key
aws configure
aws sts get-caller-identity
4. Install AWS CDK
npm install -g aws-cdk
cdk --version
5. cd to project folder
cdk bootstrap aws://ACCOUNT-NUMBER/REGION
aws sts get-caller-identity --query Account --output text
aws configure get region
cdk bootstrap aws://ACCOUNT-NUMBER/eu-central-1

API calls
1. Register
curl -X POST AWS_SERVER_URL/register -H "Content-Type: application/json" -d '{"username":"USERNAME", "password":"PASSWORD"}'
2. Login
curl -X POST AWS_SERVER_URL/login -H "Content-Type: application/json" -d '{"username":"USERNAME", "password":"PASSWORD"}'
3. Access Protected Route
curl -X GET AWS_SERVER_URL/protected -H "Content-Type: application/json" -H "Authorization: Bearer JWT_TOKEN"

Delete all infrastructure
cdk destory


Reasons
1. Quick deployment
2. No ssl overhead
3. Build in logging and troubleshooting
4. Easy to develop everywhere
5. Infrasturucture as a code
6. Easy up all backend and easy destroy


cdk init app --language go
go get

# build on windows - use wsl

# build on linux
GOOS=linux GOARCH=amd64 go build -o bootstrap

# zip package on linux
zip function.zip bootstrap

cd ..

cdk diff
cdk deploy
cdk destory

cd lambda
make build


-= TESTS =-

- users - 

curl -X POST https://tq4028wi45.execute-api.eu-central-1.amazonaws.com/prod/register -H "Content-Type: application/json" -d '{"username":"user1", "password":"password123"}'

curl -X POST https://tq4028wi45.execute-api.eu-central-1.amazonaws.com/prod/login -H "Content-Type: application/json" -d '{"username":"user1", "password":"password123"}'

curl -X GET https://tq4028wi45.execute-api.eu-central-1.amazonaws.com/prod/me -H "Content-Type: application/json" -H "Authorization: Bearer ACCESS-TOKEN"

curl -X PUT https://tq4028wi45.execute-api.eu-central-1.amazonaws.com/prod/role -H "Content-Type: application/json" -H "Authorization: Bearer ACCESS-TOKEN" -d '{"username":"user1", "newrole":"admin"}'

curl -X GET https://tq4028wi45.execute-api.eu-central-1.amazonaws.com/prod/list -H "Content-Type: application/json" -H "Authorization: Bearer ACCESS-TOKEN"

curl -X DELETE https://tq4028wi45.execute-api.eu-central-1.amazonaws.com/prod/remove?username=user111 -H "Content-Type: application/json" -H "Authorization: Bearer ACCESS-TOKEN"

- products - 

curl -X POST https://in60wqcj4h.execute-api.eu-central-1.amazonaws.com/prod/create -H "Content-Type: application/json" -H "Authorization: Bearer ACCESS-TOKEN" -d '{"name":"product1", "description":"some good product 1", "price": 101}'

curl -X GET https://in60wqcj4h.execute-api.eu-central-1.amazonaws.com/prod/list -H "Content-Type: application/json"

curl -X GET https://in60wqcj4h.execute-api.eu-central-1.amazonaws.com/prod/one?id=d396bd8f-25a2-40b9-94f2-e61942ad324a -H "Content-Type: application/json"

curl -X PUT https://in60wqcj4h.execute-api.eu-central-1.amazonaws.com/prod/update -H "Content-Type: application/json" -H "Authorization: Bearer ACCESS-TOKEN" -d '{"id": "d396bd8f-25a2-40b9-94f2-e61942ad324a", "name":"product updated", "description":"some good product updated", "price": 1000}'

curl -X DELETE https://in60wqcj4h.execute-api.eu-central-1.amazonaws.com/prod/delete?id=d396bd8f-25a2-40b9-94f2-e61942ad324a -H "Content-Type: application/json" -H "Authorization: Bearer ACCESS-TOKEN"


-= END TESTS =-

