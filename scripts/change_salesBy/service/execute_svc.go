package service

import (
	"context"
	"fmt"

	"go-scripting/pkg/logger"
	"go-scripting/repositories/dynamo_db/transactions"
	"go-scripting/repositories/firestore"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ChangeSalesByService struct {
	txV1Repository  transactions.AWSRepoFithubTransactionsV1
	txV6Repository  transactions.AWSRepoFithubTransactionsV6
	usersRepository firestore.UsersRepo
}

type ChangeSalesByInterface interface {
	ChangeSalesBy(ctx context.Context, trxId, newSalesBy string) error
}

func NewChangeSalesByService() ChangeSalesByInterface {
	txV1Repo := transactions.NewTransactionV1Repository()
	txV6Repo := transactions.NewTransactionV6Repository()
	userRepo := firestore.NewFirestoreUsersRepository()

	return &ChangeSalesByService{
		txV1Repository:  txV1Repo,
		txV6Repository:  txV6Repo,
		usersRepository: userRepo,
	}
}

type ErrorCode struct {
	Message string
}

func (e *ErrorCode) Error() string {
	return e.Message
}

// ChangeSalesBy extends the expiration date of a member's PT package.
func (salesByService *ChangeSalesByService) ChangeSalesBy(ctx context.Context, trxId, newSalesBy string) error {
	txv1Data, err := salesByService.txV1Repository.GetTransactionV1ByID(ctx, trxId)
	if err != nil {
		logger.LogInfo(fmt.Sprintf("Error processing TransactionID: %v, Error: %v", trxId, err))
		return &ErrorCode{Message: fmt.Sprintf("failed to retrieve txV1Data for TransactionID: %v", trxId)}
	}

	updateMap := map[string]interface{}{
		"salesBy":     newSalesBy,
		"salesByName": "Supriyanto",
	}
	// update v1
	err = salesByService.txV1Repository.UpdateAttributes(ctx, getStringAttribute(txv1Data, "period"), trxId, updateMap)
	if err != nil {
		return fmt.Errorf("error updating v1 data %v", trxId)
	}

	// update v6
	// check if transaction index[0] M don;t update v6
	if trxId[0] != 'M' {
		updateMapV6 := map[string]interface{}{
			"salesBy": newSalesBy,
		}
		txV6Data, err := salesByService.txV6Repository.GetTransactionV6ByID(ctx, getStringAttribute(txv1Data, "transactionNo"))
		if err != nil {
			logger.LogInfo(fmt.Sprintf("Failed to retrieve txV6Data: %v", err)) // Perbaikan: ganti log.Println dengan logger.LogInfo
			return &ErrorCode{Message: "failed to retrieve txV6Data"}           // Perbaikan: ubah pesan kesalahan untuk txV6Data
		}

		err = salesByService.txV6Repository.UpdateAttributesV6(ctx, getStringAttribute(txV6Data, "dateTransaction"), getStringAttribute(txv1Data, "transactionNo"), updateMapV6)
		if err != nil {
			return fmt.Errorf("error updating v6 data %v", trxId) // Perbaikan: ubah pesan kesalahan untuk v6
		}
	}

	// Log the successful extension.
	logger.LogInfo(fmt.Sprintf("TransactionID (%v)SalesBy changed successfully to (%v)", trxId, newSalesBy))

	return nil
}

func getStringAttribute(data map[string]types.AttributeValue, attribute string) string {
	if attr, ok := data[attribute].(*types.AttributeValueMemberS); ok {
		return attr.Value
	}
	return ""
}
