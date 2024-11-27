package service

import (
	"context"
	"go-scripting/repositories/firestore"
	"time"
)

type AdjustmentExpiredMembership struct {
	usersRepository firestore.UsersRepo
}

type AdjustmentExpiredMembershipImpl interface {
	AdjustExpiredMembership(ctx context.Context, memberId string) error
}

func NewInstanceAdjustmentExpiredMembership() AdjustmentExpiredMembershipImpl {
	return &AdjustmentExpiredMembership{
		usersRepository: firestore.NewFirestoreUsersRepository(),
	}
}

func (a *AdjustmentExpiredMembership) AdjustExpiredMembership(ctx context.Context, memberId string) error {
	member, err := a.usersRepository.QueryUserByUserAppId(memberId)
	if err != nil {
		return err
	}

	err = a.usersRepository.UpdateExpiredMembership(ctx, member.Uid, time.Now())
	if err != nil {
		return err
	}

	return nil
}
