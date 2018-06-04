package dao

import (
	"github.com/autonomousdotai/handshake-crowdfunding/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type CrowdFundingUpdateDao struct {
}

func (crowdFundingUpdateDao CrowdFundingUpdateDao) GetById(id int) (models.CrowdFundingUpdate) {
	dto := models.CrowdFundingUpdate{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingUpdateDao CrowdFundingUpdateDao) Create(dto models.CrowdFundingUpdate, tx *gorm.DB) (models.CrowdFundingUpdate, error) {
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

func (crowdFundingUpdateDao CrowdFundingUpdateDao) Update(dto models.CrowdFundingUpdate, tx *gorm.DB) (models.CrowdFundingUpdate, error) {
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

func (crowdFundingUpdateDao CrowdFundingUpdateDao) Delete(dto models.CrowdFundingUpdate, tx *gorm.DB) (models.CrowdFundingUpdate, error) {
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
