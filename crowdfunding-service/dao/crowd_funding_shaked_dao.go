package dao

import (
	"../models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type CrowdFundingShakedDao struct {
}

func (crowdFundingShakedDao CrowdFundingShakedDao) GetById(id int) (models.CrowdFundingShaked) {
	dto := models.CrowdFundingShaked{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingShakedDao CrowdFundingShakedDao) Create(dto models.CrowdFundingShaked, tx *gorm.DB) (models.CrowdFundingShaked, error) {
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

func (crowdFundingShakedDao CrowdFundingShakedDao) Update(dto models.CrowdFundingShaked, tx *gorm.DB) (models.CrowdFundingShaked, error) {
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

func (crowdFundingShakedDao CrowdFundingShakedDao) Delete(dto models.CrowdFundingShaked, tx *gorm.DB) (models.CrowdFundingShaked, error) {
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
