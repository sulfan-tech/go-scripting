package userpool

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	dynamoDB "go-scripting/configs/dynamo_db"
)

type UserPoolService struct {
	cognitoClient *cognitoidentityprovider.Client
}

func NewUserPoolService() *UserPoolService {
	awsInstance := dynamoDB.NewAWS()
	cognitoClient := awsInstance.Cognito()
	return &UserPoolService{
		cognitoClient: cognitoClient,
	}
}

func (ups *UserPoolService) GetUserByEmail(email string) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(os.Getenv("AWS_USER_POOLS_ID")),
		Username:   aws.String(email),
	}

	result, err := ups.cognitoClient.AdminGetUser(context.TODO(), input)
	if err != nil {
		fmt.Println("Error getting user:", err)
		return
	}

	fmt.Println("User information:")
	fmt.Printf("Username: %s\n", *result.Username)
	for _, attr := range result.UserAttributes {
		fmt.Printf("%s: %s\n", *attr.Name, *attr.Value)
	}
}

func (ups *UserPoolService) UpdateUserAttributes(email string, attributes map[string]string) {
	// Buat slice untuk menyimpan atribut pengguna yang akan diperbarui
	var userAttributes []types.AttributeType
	for key, value := range attributes {
		userAttributes = append(userAttributes, types.AttributeType{
			Name:  aws.String(key),
			Value: aws.String(value),
		})
	}

	input := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     aws.String(os.Getenv("AWS_USER_POOLS_ID")),
		Username:       aws.String(email),
		UserAttributes: userAttributes,
	}

	_, err := ups.cognitoClient.AdminUpdateUserAttributes(context.TODO(), input)
	if err != nil {
		fmt.Println("Error updating user attributes:", err)
		return
	}

	fmt.Println("User attributes updated successfully")
}
