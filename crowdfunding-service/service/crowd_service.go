package service

import (
	"../models"
	"../bean"
	"errors"
	"mime/multipart"
	"strings"
	"time"
	"../setting"
	"log"
	"../request_obj"
)

type CrowdService struct {
}

func (crowdService CrowdService) CreateCrowdFunding(userId int, request request_obj.CrowdFundingRequest, imageFile *multipart.File, imageFileHeader *multipart.FileHeader) (models.CrowdFunding, *bean.AppError) {
	crowdFunding := models.CrowdFunding{}
	fileNameImage := ""
	if imageFile != nil && imageFileHeader != nil {
		configuration := setting.CurrentConfig()
		uploadImageFolder := configuration.UploadImageFolder + "/crowd-funding"
		fileName := imageFileHeader.Filename
		imageExt := strings.Split(fileName, ".")[1]
		imageName := "cf-" + time.Now().Format("20060102150405")
		fileNameImage = imageName + "." + imageExt
		err := s3Service.UploadFormFile(*imageFile, uploadImageFolder, fileNameImage, imageFileHeader)
		if err != nil {
			log.Println(err)
			return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
		}
	}
	crowdFunding = models.CrowdFunding{
		Name:             request.Name,
		Description:      request.Description,
		ShortDescription: request.ShortDescription,
		CrowdDate:        request.CrowdDate,
		DeliverDate:      request.DeliverDate,
		Price:            request.Price,
		Goal:             request.Goal,
		Status:           request.Status,
	}
	crowdFunding, err := crowdFundingDao.Create(crowdFunding, nil)
	if err != nil {
		log.Println(err)
		return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}
	return crowdFunding, nil
}

func (crowdService CrowdService) GetCrowdFunding(userId int, crowdFundingId int) (models.CrowdFunding, *bean.AppError) {
	crowdFunding := crowdFundingDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFunding, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	return crowdFunding, nil
}
