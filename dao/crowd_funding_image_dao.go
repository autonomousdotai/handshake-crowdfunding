package dao

import (
	"github.com/ninjadotorg/handshake-crowdfunding/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type CrowdFundingImageDao struct {
}

func (crowdFundingImageDao CrowdFundingImageDao) GetById(id int64) (models.CrowdFundingImage) {
	dto := models.CrowdFundingImage{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingImageDao CrowdFundingImageDao) Create(dto models.CrowdFundingImage, tx *gorm.DB) (models.CrowdFundingImage, error) {
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

func (crowdFundingImageDao CrowdFundingImageDao) Update(dto models.CrowdFundingImage, tx *gorm.DB) (models.CrowdFundingImage, error) {
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

func (crowdFundingImageDao CrowdFundingImageDao) Delete(dto models.CrowdFundingImage, tx *gorm.DB) (models.CrowdFundingImage, error) {
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

func (crowdFundingImageDao CrowdFundingImageDao) GetByCrowdId(crowdFundingId int64) ([]models.CrowdFundingImage) {
	dtos := []models.CrowdFundingImage{}
	err := models.Database().Where("crowd_funding_id = ?", crowdFundingId).Find(&dtos).Error
	if err != nil {
		log.Print(err)
	}
	return dtos
}
