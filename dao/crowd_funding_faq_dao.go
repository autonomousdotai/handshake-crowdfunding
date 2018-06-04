package dao

import (
	"github.com/ninjadotorg/handshake-crowdfunding/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type CrowdFundingFaqDao struct {
}

func (crowdFundingFaqDao CrowdFundingFaqDao) GetById(id int) (models.CrowdFundingFaq) {
	dto := models.CrowdFundingFaq{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingFaqDao CrowdFundingFaqDao) Create(dto models.CrowdFundingFaq, tx *gorm.DB) (models.CrowdFundingFaq, error) {
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

func (crowdFundingFaqDao CrowdFundingFaqDao) Update(dto models.CrowdFundingFaq, tx *gorm.DB) (models.CrowdFundingFaq, error) {
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

func (crowdFundingFaqDao CrowdFundingFaqDao) Delete(dto models.CrowdFundingFaq, tx *gorm.DB) (models.CrowdFundingFaq, error) {
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
