package dao

import (
	"github.com/ninjadotorg/handshake-crowdfunding/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/ninjadotorg/handshake-crowdfunding/bean"
)

type CrowdFundingFaqDao struct {
}

func (crowdFundingFaqDao CrowdFundingFaqDao) GetById(id int64) (models.CrowdFundingFaq) {
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

func (crowdFundingFaqDao CrowdFundingFaqDao) GetAllBy(userId int64, crowdFundingId int64, pagination *bean.Pagination) (*bean.Pagination, error) {
	dtos := []models.CrowdFundingFaq{}
	db := models.Database()
	if pagination != nil {
		db = db.Limit(pagination.PageSize)
		db = db.Offset(pagination.PageSize * (pagination.Page - 1))
	}
	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}
	if crowdFundingId > 0 {
		db = db.Where("crowd_funding_id = ?", crowdFundingId)
	}
	err := db.Order("prioriry asc, date_created desc").Find(&dtos).Error
	if err != nil {
		log.Print(err)
		return pagination, err
	}
	pagination.Items = dtos
	total := 0
	if pagination.Page == 1 && len(dtos) < pagination.PageSize {
		total = len(dtos)
	} else {
		err := db.Find(&dtos).Count(&total).Error
		if err != nil {
			log.Print(err)
			return pagination, err
		}
	}
	pagination.Total = total
	return pagination, nil
}
