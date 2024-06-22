package transactions

import (
	"context"
	"errors"
	"fmt"
	dynamoDB "go-scripting/configs/dynamo_db"
	"go-scripting/pkg/constant"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type dynamoDBRepo struct {
	dydb *dynamodb.Client
}

type AWSRepoFithubTransactionsV6 interface {
	GetTransactionV6ByID(ctx context.Context, transId string) (map[string]types.AttributeValue, error)
	UpdateFieldsV6(ctx context.Context, transId, newLocationTransaction, dateTransaction, newSalesBy string) error
	UpdateAttributesV6(ctx context.Context, dateTransaction, transId string, updateMap map[string]interface{}) error
	InsertDocumentV6(ctx context.Context, item map[string]types.AttributeValue) error
	DeleteTransactionV6(ctx context.Context, dateTransaction, transId string) error
}

func NewTransactionV6Repository() AWSRepoFithubTransactionsV6 {
	awsCfg := dynamoDB.NewAWS()
	client := awsCfg.DynamoDB()

	return &dynamoDBRepo{
		dydb: client,
	}
}

func (d *dynamoDBRepo) GetTransactionV6ByID(ctx context.Context, transId string) (map[string]types.AttributeValue, error) {

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(constant.TABLE_NAME_CS_TRANSACTION_V6),
		IndexName:              aws.String("transId-index"),
		KeyConditionExpression: aws.String("transId = :transId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":transId": &types.AttributeValueMemberS{
				Value: transId,
			},
		},
	}

	// Create a QueryPaginator
	paginator := dynamodb.NewQueryPaginator(d.dydb, queryInput)

	// Iterate through the pages
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if len(page.Items) > 0 {
			return page.Items[0], nil
		}
	}

	// No items found
	return nil, errors.New("transaction not found")
}

func (d *dynamoDBRepo) UpdateFieldsV6(ctx context.Context, transId, dateTransaction, newLocationTransaction, newSalesBy string) error {

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(constant.TABLE_NAME_CS_TRANSACTION_V6),
		Key: map[string]types.AttributeValue{
			"dateTransaction": &types.AttributeValueMemberS{
				Value: dateTransaction,
			},
			"transId": &types.AttributeValueMemberS{
				Value: transId,
			},
		},
		UpdateExpression: aws.String("SET #salesBy =:newSalesBy, #locationOn =:newLocationTransaction REMOVE #requestSalesBy, #isRejected, #isRequested"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newSalesBy":             &types.AttributeValueMemberS{Value: newSalesBy},
			":newLocationTransaction": &types.AttributeValueMemberS{Value: newLocationTransaction},
		},
		ExpressionAttributeNames: map[string]string{
			"#salesBy":        "salesBy",
			"#locationOn":     "locationOn",
			"#requestSalesBy": "requestSalesBy",
			"#isRejected":     "isRejected",
			"#isRequested":    "isRequested",
		},
	}

	_, err := d.dydb.UpdateItem(ctx, updateItemInput)
	if err != nil {
		fmt.Println("Error updating item:", err)
		return err
	}

	return nil
}

func (d *dynamoDBRepo) InsertDocumentV6(ctx context.Context, item map[string]types.AttributeValue) error {
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(constant.TABLE_NAME_CS_TRANSACTION_V6),
		Item:      item,
	}

	_, err := d.dydb.PutItem(ctx, putInput)
	if err != nil {
		fmt.Println("Error inserting item:", err)
		return err
	}

	return nil
}

func (d *dynamoDBRepo) DeleteTransactionV6(ctx context.Context, dateTransaction, transId string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(constant.TABLE_NAME_CS_TRANSACTION_V6),
		Key: map[string]types.AttributeValue{
			"dateTransaction": &types.AttributeValueMemberS{
				Value: dateTransaction,
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

func (d *dynamoDBRepo) UpdateAttributesV6(ctx context.Context, dateTransaction, transId string, updateMap map[string]interface{}) error {
	updateExpression, expressionAttributeValues, expressionAttributeNames := buildUpdateExpression(updateMap)

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(constant.TABLE_NAME_CS_TRANSACTION_V6),
		Key:                       map[string]types.AttributeValue{"dateTransaction": &types.AttributeValueMemberS{Value: dateTransaction}, "transId": &types.AttributeValueMemberS{Value: transId}},
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
