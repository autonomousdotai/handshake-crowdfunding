package dao

import (
	"../models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type CrowdFundingDao struct {
}

func (crowdFundingDao CrowdFundingDao) GetById(id int64) (models.CrowdFunding) {
	dto := models.CrowdFunding{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingDao CrowdFundingDao) Create(dto models.CrowdFunding, tx *gorm.DB) (models.CrowdFunding, error) {
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

func (crowdFundingDao CrowdFundingDao) Update(dto models.CrowdFunding, tx *gorm.DB) (models.CrowdFunding, error) {
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

func (crowdFundingDao CrowdFundingDao) Delete(dto models.CrowdFunding, tx *gorm.DB) (models.CrowdFunding, error) {
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
