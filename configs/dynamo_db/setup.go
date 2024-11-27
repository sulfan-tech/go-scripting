package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type AwsConfig struct {
	cfg aws.Config
}

func NewAWS() *AwsConfig {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	return &AwsConfig{
		cfg: cfg,
	}
}

func (ac *AwsConfig) DynamoDB() *dynamodb.Client {
	return dynamodb.NewFromConfig(ac.cfg)
}

func (ac *AwsConfig) Cognito() *cognitoidentityprovider.Client {
	return cognitoidentityprovider.NewFromConfig(ac.cfg)
}
