package mysql

import (
	sqlDB "go-scripting/configs/mySql"
	"go-scripting/entities"

	"github.com/jinzhu/gorm"
)

type EBARepository interface {
	CreateUserEBA(param entities.EBA) error
	BypassEBA(param entities.BYPASS) error
}

type EBARepositoryImpl struct {
	db *gorm.DB
}

func NewEBARepository() (EBARepository, error) {
	db, err := sqlDB.SetupDatabase()
	if err != nil {
		return nil, err
	}

	return &EBARepositoryImpl{db}, nil
}

func (repo *EBARepositoryImpl) UpdateExpire(uid string) error {

	return nil
}

func (repo *EBARepositoryImpl) CreateUserEBA(param entities.EBA) error {

	qr := repo.db.Table("early_access_hoard").Create(&param)

	// rowrAffected := qr.RowsAffected

	if err := qr.Error; err != nil {
		return err
	}

	// fmt.Println(rowrAffected)
	// if err := repo.db.Table("early_access_hoard").Create(&param).Error; err != nil {
	// 	return err
	// }
	return nil
}

func (repo *EBARepositoryImpl) BypassEBA(param entities.BYPASS) error {
	// if user is already exist just update not created
	if err := repo.db.Table("early_access_bypass").Create(&param).Error; err != nil {
		return err
	}
	return nil
}
