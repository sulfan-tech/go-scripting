package service

import (
	"context"
	"fmt"
	"go-scripting/entities"
	repoMySql "go-scripting/repositories/mySql"
	"log"
)

type EBAService struct {
	ebaRepo repoMySql.EBARepository
}

type EBAInterface interface {
	InsertUserEBA(ctx context.Context, param RequestUserEBA) error
}

func NewUserEBAService() EBAInterface {
	ebaRepo, err := repoMySql.NewEBARepository()
	if err != nil {
		log.Fatal("error creating ebaRepo repository")
	}

	return &EBAService{
		ebaRepo: ebaRepo,
	}
}

type RequestUserEBA struct {
	InsertEBA entities.EBA
	BypassEBA entities.BYPASS
}

type ErrorCode struct {
	Message string
}

func (e *ErrorCode) Error() string {
	return e.Message
}

func (r *EBAService) InsertUserEBA(ctx context.Context, param RequestUserEBA) error {

	err := r.ebaRepo.CreateUserEBA(param.InsertEBA)
	if err != nil {
		return &ErrorCode{Message: fmt.Sprintf("error create userEBA: %v", err)}
	}

	err = r.ebaRepo.BypassEBA(param.BypassEBA)
	if err != nil {
		return &ErrorCode{Message: fmt.Sprintf("error bypass userEBA: %v", err)}
	}

	return nil
}
