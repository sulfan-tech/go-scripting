package transactions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	dynamoDB "go-scripting/configs/dynamo_db"

	"go-scripting/pkg/constant"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoDBRepository struct {
	dydb *dynamodb.Client
}

type AWSRepoFithubTransactionsV1 interface {
	GetTransactionV1ByID(ctx context.Context, transId string) (map[string]types.AttributeValue, error)
	UpdateFieldsV1(ctx context.Context, transId, period, newSalesBy, newLocationTransaction string) error
	UpdateTimeFieldsV1(ctx context.Context, transId, period, transDate, transDateJakarta, transDateUTC, pushDate string) error
	UpdateAttributes(ctx context.Context, period, transId string, updateMap map[string]interface{}) error
	InsertDocumentV1(ctx context.Context, item map[string]types.AttributeValue) error
	DeleteTransactionV1(ctx context.Context, period, transId string) error
}

func NewTransactionV1Repository() AWSRepoFithubTransactionsV1 {
	awsCfg := dynamoDB.NewAWS()
	client := awsCfg.DynamoDB()

	return &dynamoDBRepository{
		dydb: client,
	}
}

func (d *dynamoDBRepository) GetTransactionV1ByID(ctx context.Context, transId string) (map[string]types.AttributeValue, error) {

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(constant.TABLE_NAME_TRANSACTION_V1),
		IndexName:              aws.String("byTransId_index"),
		KeyConditionExpression: aws.String("transId = :transId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":transId": &types.AttributeValueMemberS{Value: transId},
		},
	}

	p := dynamodb.NewQueryPaginator(d.dydb, queryInput)

	for p.HasMorePages() {
		output, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if len(output.Items) > 0 {
			return output.Items[0], nil
		}
	}

	// No items found
	return nil, nil
}

func (d *dynamoDBRepository) UpdateFieldsV1(ctx context.Context, transId, period, newSalesBy, newLocationTransaction string) error {

	// key + salesby + locationTransaction
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(constant.TABLE_NAME_TRANSACTION_V1),
		Key: map[string]types.AttributeValue{
			"transId": &types.AttributeValueMemberS{
				Value: transId,
			},
			"period": &types.AttributeValueMemberS{
				Value: period,
			},
		},
		UpdateExpression: aws.String("SET #salesBy =:newSalesBy, #locationTransaction =:newLocationTransaction REMOVE #requestSalesBy, #isRejected, #isRequested"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newSalesBy":             &types.AttributeValueMemberS{Value: newSalesBy},
			":newLocationTransaction": &types.AttributeValueMemberS{Value: newLocationTransaction},
		},
		ExpressionAttributeNames: map[string]string{
			"#salesBy":             "salesBy",
			"#locationTransaction": "locationTransaction",
			"#requestSalesBy":      "requestSalesBy",
			"#isRejected":          "isRejected",
			"#isRequested":         "isRequested",
		},
	}

	_, err := d.dydb.UpdateItem(ctx, updateItemInput)
	if err != nil {
		fmt.Println("Error updating item:", err)
		return err
	}

	return nil
}

func (d *dynamoDBRepository) UpdateTimeFieldsV1(ctx context.Context, transId, period, transDate, transDateJakarta, transDateUTC, pushDate string) error {
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(constant.TABLE_NAME_TRANSACTION_V1),
		Key: map[string]types.AttributeValue{
			"transId": &types.AttributeValueMemberS{
				Value: transId,
			},
			"period": &types.AttributeValueMemberS{
				Value: period,
			},
		},
		UpdateExpression: aws.String("SET #transDate = :transDate, #transDate_jakarta = :transDateJakarta, #transDate_utc = :transDateUTC, #pushDate = :pushDate"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":transDate":        &types.AttributeValueMemberS{Value: transDate},
			":transDateJakarta": &types.AttributeValueMemberS{Value: transDateJakarta},
			":transDateUTC":     &types.AttributeValueMemberS{Value: transDateUTC},
			":pushDate":         &types.AttributeValueMemberS{Value: pushDate},
		},
		ExpressionAttributeNames: map[string]string{
			"#transDate":         "transDate",
			"#transDate_jakarta": "transDate_jakarta",
			"#transDate_utc":     "transDate_utc",
			"#pushDate":          "pushDate",
		},
	}

	_, err := d.dydb.UpdateItem(ctx, updateItemInput)
	if err != nil {
		fmt.Println("Error updating item:", err)
		return err
	}

	return nil
}

func (d *dynamoDBRepository) InsertDocumentV1(ctx context.Context, item map[string]types.AttributeValue) error {
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(constant.TABLE_NAME_TRANSACTION_V1),
		Item:      item,
	}

	_, err := d.dydb.PutItem(ctx, putInput)
	if err != nil {
		fmt.Println("Error inserting item:", err)
		return err
	}

	return nil
}

func (d *dynamoDBRepository) DeleteTransactionV1(ctx context.Context, period, transId string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(constant.TABLE_NAME_TRANSACTION_V1),
		Key: map[string]types.AttributeValue{
			"period": &types.AttributeValueMemberS{
				Value: period,
			},
			"transId": &types.AttributeValueMemberS{
				Value: transId,
			},
		},
	}

	_, err := d.dydb.DeleteItem(ctx, deleteInput)
	if err != nil {
		fmt.Println("Error deleting item:", err)
		return err
	}

	return nil
}

func (d *dynamoDBRepository) UpdateAttributes(ctx context.Context, period, transId string, updateMap map[string]interface{}) error {
	updateExpression, expressionAttributeValues, expressionAttributeNames := buildUpdateExpression(updateMap)

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(constant.TABLE_NAME_TRANSACTION_V1),
		Key:                       map[string]types.AttributeValue{"period": &types.AttributeValueMemberS{Value: period}, "transId": &types.AttributeValueMemberS{Value: transId}},
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
