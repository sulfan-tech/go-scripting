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
	AddFreeDays(ctx context.Context, memberId, daysToAdd, logs string) error
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
func (f *FreeDaysService) AddFreeDays(ctx context.Context, memberId, daysToAdd, logs string) error {

	user, err := f.usersRepository.QueryUserByUserAppId(memberId)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error querying user by memberId (%s): %v", memberId, err))
		return fmt.Errorf("error querying user by memberId (%s): %w", memberId, err)
	}

	if user == nil {
		logger.LogInfo(fmt.Sprintf("User with memberId (%s) not found", memberId))
		return fmt.Errorf("user with memberId (%s) not found", memberId)
	}

	days, err := strconv.Atoi(daysToAdd)
	if err != nil {
		return fmt.Errorf("failed to parse daysToAdd: %w", err)
	}

	newExpiredMembership := user.ExpiredMembership.AddDate(0, 0, days)
	if err := f.usersRepository.UpdateExpiredMembership(ctx, user.Uid, newExpiredMembership); err != nil {
		logger.LogError(fmt.Sprintf("Error updating expired membership for user with memberId (%s): %v", memberId, err))
		return fmt.Errorf("error updating expired membership for user with memberId (%s): %w", memberId, err)
	}

	time.Sleep(7 * time.Second)

	var logMemberships []entities.UserLogsMembership
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		logMemberships, err = f.usersRepository.GetLatestLogMembership(ctx, user.Uid)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error getting latest log membership for user with memberId (%s): %v", memberId, err))
			return fmt.Errorf("error getting latest log membership for user with memberId (%s): %w", memberId, err)
		}

		if len(logMemberships) > 0 {
			break
		}

		if i < maxRetries-1 {
			time.Sleep(5 * time.Second)
		}
	}

	if len(logMemberships) == 0 {
		logger.LogError(fmt.Sprintf("No logsMembership found in the latest log after multiple attempts for user with memberId (%s)", memberId))
		return fmt.Errorf("no logsMembership found in the latest log after multiple attempts for user with memberId (%s)", memberId)
	}

	logsMembership := logs
	if logsMembership == "" {
		logsMembership = "Adjustment ExpiredDate (Ticketing)"
	}

	for _, log := range logMemberships {
		if log.NewDate.Equal(newExpiredMembership) {
			if err := f.usersRepository.UpdateLogMembershipAndTypeChange(ctx, user.Uid, logsMembership); err != nil {
				logger.LogError(fmt.Sprintf("Error updating log membership for user with memberId (%s): %v", memberId, err))
				return fmt.Errorf("error updating log membership for user with memberId (%s): %w", memberId, err)
			}
			break
		}
	}

	return nil
}
