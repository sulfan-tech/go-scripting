package service

import (
	"context"
	"fmt"
	"go-scripting/entities"
	"go-scripting/pkg/logger"
	"go-scripting/repositories/firestore"
	"log"
	"time"
)

type UserDummyService struct {
	usersRepository firestore.UsersRepo
}

type UserDummyImpl interface {
	CreateUserDummyFS(ctx context.Context, count int) error
}

func NewInstanceDummyUserService() UserDummyImpl {
	return &UserDummyService{
		usersRepository: firestore.NewFirestoreUsersRepository(),
	}
}

func (u *UserDummyService) CreateUserDummyFS(ctx context.Context, count int) error {
	for i := 0; i < count; i++ {
		user := entities.User{
			CreatedDate:       time.Date(2022, time.December, 6, 13, 35, 53, 0, time.UTC),
			CRMId:             "dummy",
			DateOfBirth:       time.Date(2022, time.December, 6, 7, 0, 0, 0, time.UTC),
			Email:             fmt.Sprintf("dummy%v@gmail.com", i),
			ExpiredMembership: time.Now().AddDate(0, 1, 0),
			Gender:            "Male",
			Name:              fmt.Sprintf("dummy user %v", i),
			NameLower:         fmt.Sprintf("dummy user %v", i),
			PackageName:       "FITHUB Corporate 3 Months - Rp 930,000.00",
			PackagePrice:      930000,
			Partners:          "-",
			Phone:             fmt.Sprintf("+628524628%v", i),
			PhotoIdentify:     nil,
			PhotoUrl:          "https://s.gravatar.com/avatar/1968f1b32d23b4467de1fba99a54b037?s=200&d=identicon",
			PromoLocation:     "FIT HUB ARTERI PONDOK INDAH",
			RefCode:           fmt.Sprintf("refCode%d", i),
			RegUserAppId:      "increment start1",
			SalesName:         "",
			Source:            "",
			StartDate:         time.Now(),
			ThumbLevel:        1,
			TransDate:         "",
			TransId:           "",
			TypeUser:          "NEW",
			Uid:               fmt.Sprintf("UIDD%d", i),
			UserAppId:         fmt.Sprintf("TDUMMY%d", i),
		}

		err := u.usersRepository.CreateUser(ctx, &user)
		if err != nil {
			log.Printf("Failed to create user %d: %v\n", i, err)
			return err
		}

		logger.LogInfo(user.UserAppId)
		fmt.Println("Created user", user.UserAppId)
	}

	return nil
}
