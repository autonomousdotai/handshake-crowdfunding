package dao

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-crowdfunding/bean"
	"github.com/ninjadotorg/handshake-crowdfunding/models"
)

type CrowdFundingPostDao struct {
}

func (crowdFundingPostDao CrowdFundingPostDao) GetById(id int64) models.CrowdFundingPost {
	dto := models.CrowdFundingPost{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingPostDao CrowdFundingPostDao) Create(dto models.CrowdFundingPost, tx *gorm.DB) (models.CrowdFundingPost, error) {
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

func (crowdFundingPostDao CrowdFundingPostDao) Update(dto models.CrowdFundingPost, tx *gorm.DB) (models.CrowdFundingPost, error) {
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

func (crowdFundingPostDao CrowdFundingPostDao) Delete(dto models.CrowdFundingPost, tx *gorm.DB) (models.CrowdFundingPost, error) {
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

func (crowdFundingPostDao CrowdFundingPostDao) GetAllBy(userId int64, crowdFundingId int64, pagination *bean.Pagination) (*bean.Pagination, error) {
	dtos := []models.CrowdFundingPost{}
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
	err := db.Order("date_created desc").Find(&dtos).Error
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
