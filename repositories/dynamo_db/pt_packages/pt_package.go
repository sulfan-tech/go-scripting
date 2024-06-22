package ptpackages

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	dynamoDB "go-scripting/configs/dynamo_db"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type dynamoDBRepository struct {
	dydb *dynamodb.Client
}

type PTPackagesRepo interface {
	UpdatePTPackages(ctx context.Context, ptId string, voucherId string, updateMap map[string]interface{}) error
}

func NewPTPackageRepository() PTPackagesRepo {
	awsCfg := dynamoDB.NewAWS()
	client := awsCfg.DynamoDB()

	return &dynamoDBRepository{
		dydb: client,
	}
}

func (d *dynamoDBRepository) UpdatePTPackages(ctx context.Context, ptId string, voucherId string, updateMap map[string]interface{}) error {
	updateExpression, expressionAttributeValues, expressionAttributeNames := buildUpdateExpression(updateMap)

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("fithub_pt_packages"),
		Key:                       map[string]types.AttributeValue{"ptId": &types.AttributeValueMemberS{Value: ptId}, "voucherId": &types.AttributeValueMemberS{Value: voucherId}},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
	}

	_, err := d.dydb.UpdateItem(ctx, updateItemInput)
	if err != nil {
		fmt.Println("Error updating item:", err)
		return err
	}

	return nil
}

func buildUpdateExpression(updateMap map[string]interface{}) (string, map[string]types.AttributeValue, map[string]string) {
	var updateExpression string
	expressionAttributeValues := make(map[string]types.AttributeValue)
	expressionAttributeNames := make(map[string]string)

	// Separate variables for SET and placeholders
	var setExpressions []string

	for key, value := range updateMap {
		attributeValue := valueToStringAttributeValue(value)
		if attributeValue != nil {
			placeholder := fmt.Sprintf(":%s", key)
			attributeName := fmt.Sprintf("#%s", key)

			setExpressions = append(setExpressions, fmt.Sprintf("%s = %s", attributeName, placeholder))
			expressionAttributeValues[placeholder] = attributeValue
			expressionAttributeNames[attributeName] = key
		}
	}

	// Combine all SET expressions
	if len(setExpressions) > 0 {
		updateExpression = "SET " + strings.Join(setExpressions, ", ")
	}

	return updateExpression, expressionAttributeValues, expressionAttributeNames
}

func valueToStringAttributeValue(value interface{}) types.AttributeValue {
	switch v := value.(type) {
	case int:
		return &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", v)}
	case float64:
		return &types.AttributeValueMemberN{Value: strconv.FormatFloat(v, 'f', -1, 64)}
	case string:
		return &types.AttributeValueMemberS{Value: v}
	default:
		// Handle other data types as needed
		return nil
	}
}
