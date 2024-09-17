package common

import (
	"github.com/google/uuid"
)

const ProductTableName = "JITestDemoProductTable"
const TokenSecret = "very-strong-secret"
const RoleUser = "user"
const RoleAdmin = "admin"

func GenerateStrignID() string {
	id := uuid.New()
	return id.String()
}
