package service

import (
	"context"
	"fmt"
	"go-scripting/entities"
	"go-scripting/pkg/logger"
	"go-scripting/repositories/firestore"
	"strconv"
	"time"
)

type FreeDaysService struct {
	usersRepository firestore.UsersRepo
}

type FreeDaysInterface interface {
	AddFreeDays(ctx context.Context, ExecutionId, memberId, daysToAdd, logs string) (bool, bool, error)
}

type ErrorCode struct {
	Message string
}

func (e *ErrorCode) Error() string {
	return e.Message
}

func NewFreeDaysService() FreeDaysInterface {
	return &FreeDaysService{
		usersRepository: firestore.NewFirestoreUsersRepository(),
	}
}

// AddFreeDays adds free days to a user's membership and updates logs accordingly
func (f *FreeDaysService) AddFreeDays(ctx context.Context, executionId, memberId, daysToAdd, logs string) (bool, bool, error) {
	isSuccessUpdateExpiredMembership, isSuccessTypeChange := false, false

	// Query user by memberId
	user, err := f.usersRepository.QueryUserByUserAppId(memberId)
	if err != nil {
		return false, false, fmt.Errorf("error querying memberId (%s): %w", memberId, err)
	}
	if user == nil {
		return false, false, fmt.Errorf("memberId (%s) not found", memberId)
	}

	days, err := strconv.Atoi(daysToAdd)
	if err != nil {
		return false, false, fmt.Errorf("failed to parse daysToAdd: %w", err)
	}

	// the current date at the beginning of the day in UTC+7 timezone
	currentDate := time.Now().UTC().Add(-24 * time.Hour).Add(7 * time.Hour)

	if user.MembershipStatus == 1 && user.ExpiredUpdate.After(currentDate) {
		return false, false, fmt.Errorf("memberId (%s) has a transfer provider status", memberId)
	}

	// Check if executionId already exists in the database
	existingLog, errExistingLog := f.usersRepository.GetLogMembershipByExecutionId(ctx, user.Uid, executionId)
	if errExistingLog != nil {
		return false, false, fmt.Errorf("error checking existing log membership: %w", errExistingLog)
	}
	if existingLog != nil && existingLog.ExecutionId == executionId {
		return false, false, fmt.Errorf("user already updated with execution id: %s", executionId)
	}

	// Update expiredMembership with skip trigger function
	newExpiredMembership := user.ExpiredMembership.AddDate(0, 0, days)
	if err := f.usersRepository.UpdateExpiredMembership(ctx, user.Uid, newExpiredMembership); err != nil {
		logger.LogError(fmt.Sprintf("Error updating expired membership for user with memberId (%s): %v", memberId, err))
		return false, false, fmt.Errorf("error updating expired membership for user with memberId (%s): %w", memberId, err)
	}
	isSuccessUpdateExpiredMembership = true

	// Create logs membership
	logMembership := entities.UserLogsMembership{
		ExecutionId:  executionId,
		DateTime:     time.Now(),
		NewDate:      newExpiredMembership,
		PreviousDate: user.ExpiredMembership,
		TypeChange:   logs,
	}
	err = f.usersRepository.CreateLogMembership(ctx, user.Uid, logMembership)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error creating logs membership for user with memberId (%s): %v", memberId, err))
		isSuccessTypeChange = false
		return false, false, fmt.Errorf("error creating logs membership for user with memberId (%s): %w", memberId, err)
	}
	isSuccessTypeChange = true

	return isSuccessUpdateExpiredMembership, isSuccessTypeChange, nil
}
