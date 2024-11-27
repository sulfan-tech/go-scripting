package main

import (
	userpool "go-scripting/repositories/dynamo_db/user_pool"
)

func main() {
	userPoolService := userpool.NewUserPoolService()

	userPoolService.GetUserByEmail("ninofinance@finance.id")

	// attributes := map[string]string{
	// 	"custom:attribute": "new value",
	// 	"email":            "newemail@example.com",
	// }
	// userPoolService.UpdateUserAttributes("email@example.com", attributes)
}
