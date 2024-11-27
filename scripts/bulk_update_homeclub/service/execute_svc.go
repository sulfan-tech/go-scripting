package service

import (
	"context"
	"fmt"

	"go-scripting/pkg/logger"
	"go-scripting/repositories/firestore"
)

type ChangeHomeClub struct {
	usersRepository firestore.UsersRepo
}

type ChangeHomeClubInterface interface {
	BulkChangeHomeClub(ctx context.Context, memberId, newHomeClub string) error
}

func NewChangeHomeClubService() ChangeHomeClubInterface {
	userRepo := firestore.NewFirestoreUsersRepository()

	return &ChangeHomeClub{
		usersRepository: userRepo,
	}
}

func (c *ChangeHomeClub) BulkChangeHomeClub(ctx context.Context, memberId, newHomeClub string) error {
	member, err := c.usersRepository.QueryUserByUserAppId(memberId)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error querying member by memberId (%s): %v", memberId, err))
		return fmt.Errorf("error querying member by memberId (%s): %v", memberId, err)
	}

	if member == nil {
		logger.LogInfo(fmt.Sprintf("Member with memberId (%s) not found", memberId))
		return fmt.Errorf("member with memberId (%s) not found", memberId)
	}

	updateData := map[string]interface{}{
		"promoLocation":    newHomeClub,
		"oldPromoLocation": member.PromoLocation,
	}

	err = c.usersRepository.UpdateFieldMembership(ctx, member.Uid, updateData)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error updating homeclub for member with memberId (%s): %v", memberId, err))
		return fmt.Errorf("error updating homeclub for member with memberId (%s): %v", memberId, err)
	}

	return nil
}
