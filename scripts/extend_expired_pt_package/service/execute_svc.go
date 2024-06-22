package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-scripting/pkg/constant"
	"go-scripting/pkg/logger"
	"go-scripting/repositories/dynamo_db/transactions"
	"go-scripting/repositories/firestore"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ExtendExpiredPTService struct {
	txV1Repository  transactions.AWSRepoFithubTransactionsV1
	txV6Repository  transactions.AWSRepoFithubTransactionsV6
	usersRepository firestore.UsersRepo
}

type ExtendExpiredPTInterface interface {
	ExtendExpiredPT(ctx context.Context, param RequestExtendMemberPTPackage) error
}

func NewExtendExpiredPTService() ExtendExpiredPTInterface {
	txV1Repo := transactions.NewTransactionV1Repository()
	txV6Repo := transactions.NewTransactionV6Repository()
	userRepo := firestore.NewFirestoreUsersRepository()

	return &ExtendExpiredPTService{
		txV1Repository:  txV1Repo,
		txV6Repository:  txV6Repo,
		usersRepository: userRepo,
	}
}

type RequestExtendMemberPTPackage struct {
	TransactionID     string
	MemberPhoneNumber string
	NewExpiredDate    string
	Type              string
	ExtendStatus      string
}

type ErrorCode struct {
	Message string
}

func (e *ErrorCode) Error() string {
	return e.Message
}

// ExtendExpiredPT extends the expiration date of a member's PT package.
func (ep *ExtendExpiredPTService) ExtendExpiredPT(ctx context.Context, param RequestExtendMemberPTPackage) error {
	newExpiredTime, err := time.Parse(constant.DateFormat, param.NewExpiredDate)
	if err != nil {
		return &ErrorCode{Message: fmt.Sprintf("error parsing newDate: %v", err)}
	}

	txv1Data, err := ep.txV1Repository.GetTransactionV1ByID(ctx, param.TransactionID)
	if err != nil {
		log.Println("Error retrieving txV1Data:", err)
		return &ErrorCode{Message: "failed to retrieve txV1Data"}
	}

	transactionNoValue, ok := txv1Data["transactionNo"]
	if !ok {
		return &ErrorCode{Message: "transactionNo not found"}
	}

	transactionNo, ok := transactionNoValue.(*types.AttributeValueMemberS)
	if !ok {
		return &ErrorCode{Message: "transactionNo is not a string"}
	}

	transactionNoStr := transactionNo.Value

	member, err := ep.usersRepository.QueryUserByPhone(param.MemberPhoneNumber)
	if err != nil {
		return &ErrorCode{Message: "Member is not found"}
	}

	// Update user's PT package expiration date.
	updateFields := map[string]interface{}{
		"expiredDate": newExpiredTime,
	}

	err = ep.usersRepository.UpdateUserPTPackageFields(ctx, member.Uid, transactionNoStr, updateFields)
	if err != nil {
		return &ErrorCode{Message: "Error updating pt packages"}
	}

	// Retrieve member's PT package details.
	memberPT, err := ep.usersRepository.GetPTPackagesByUID(ctx, member.Uid, transactionNoStr)
	if err != nil {
		return &ErrorCode{Message: "Member PT Package is not found"}
	}

	oldExpiredPT := memberPT.ExpiredDate
	existingPackageName := memberPT.PackageName

	// Log the successful extension.
	logger.LogInfo(fmt.Sprintf("Extend Expired PT Package [%v] - Package Name [%v] from Member Phone (%v) successfully extended to (%v)", oldExpiredPT, existingPackageName, param.MemberPhoneNumber, newExpiredTime))

	return nil
}
