package dao

import (
	"github.com/ninjadotorg/handshake-crowdfunding/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type CrowdFundingShakeDao struct {
}

func (crowdFundingShakeDao CrowdFundingShakeDao) GetById(id int64) (models.CrowdFundingShake) {
	dto := models.CrowdFundingShake{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingShakedDao CrowdFundingShakeDao) Create(dto models.CrowdFundingShake, tx *gorm.DB) (models.CrowdFundingShake, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateCreated = time.Now()
	dto.DateModified = dto.DateCreated
	err := tx.Create(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (crowdFundingShakedDao CrowdFundingShakeDao) Update(dto models.CrowdFundingShake, tx *gorm.DB) (models.CrowdFundingShake, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateModified = time.Now()
	err := tx.Save(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (crowdFundingShakedDao CrowdFundingShakeDao) Delete(dto models.CrowdFundingShake, tx *gorm.DB) (models.CrowdFundingShake, error) {
	if tx == nil {
		tx = models.Database()
	}
	err := tx.Delete(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (crowdFundingShakeDao CrowdFundingShakeDao) GetAllByBackerStatus(crowdFundingId int64, userId int64, status int) ([]models.CrowdFundingShake) {
	dtos := []models.CrowdFundingShake{}
	db := models.Database()
	if crowdFundingId > 0 {
		db = db.Where("crowd_funding_id = ?", crowdFundingId)
	}
	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	err := db.Find(&dtos).Error
	if err != nil {
		log.Print(err)
	}
	return dtos
}
