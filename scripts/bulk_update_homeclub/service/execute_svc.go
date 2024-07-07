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
	BulkChangeHomeClub(ctx context.Context, phone, newHomeClub string) error
}

func NewChangeHomeClubService() ChangeHomeClubInterface {
	userRepo := firestore.NewFirestoreUsersRepository()

	return &ChangeHomeClub{
		usersRepository: userRepo,
	}
}

func (c *ChangeHomeClub) BulkChangeHomeClub(ctx context.Context, phone, newHomeClub string) error {
	member, err := c.usersRepository.QueryUserByPhone(phone)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error querying member by phone (%s): %v", phone, err))
		return fmt.Errorf("error querying member by phone (%s): %v", phone, err)
	}

	if member == nil {
		logger.LogInfo(fmt.Sprintf("Member with phone (%s) not found", phone))
		return fmt.Errorf("member with phone (%s) not found", phone)
	}

	updateData := map[string]interface{}{
		"promoLocation":    newHomeClub,
		"oldPromoLocation": member.PromoLocation,
	}

	err = c.usersRepository.UpdateFieldMembership(ctx, member.Uid, updateData)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error updating homeclub for member with phone (%s): %v", phone, err))
		return fmt.Errorf("error updating homeclub for member with phone (%s): %v", phone, err)
	}

	return nil
}
