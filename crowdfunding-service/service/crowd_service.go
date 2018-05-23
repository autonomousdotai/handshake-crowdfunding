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
	"../utils"
	"github.com/jinzhu/gorm"
)

type CrowdService struct {
}

func (crowdService CrowdService) CreateTx(userId int64, address string, hash string, refType string, refId int64, tx *gorm.DB) (models.EthTx, *bean.AppError) {
	ethTx := models.EthTx{}
	ethTx.UserId = userId
	ethTx.FromAddress = address
	ethTx.Hash = hash
	ethTx.RefType = refType
	ethTx.RefId = refId
	ethTx.Status = 0
	ethTx, err := ethTxDao.Create(ethTx, tx)
	if err != nil {
		log.Println(err)
		return ethTx, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}
	return ethTx, nil
}

func (crowdService CrowdService) CreateCrowdFunding(userId int64, request request_obj.CrowdFundingRequest, imageFile *multipart.File, imageFileHeader *multipart.FileHeader) (models.CrowdFunding, *bean.AppError) {
	crowdFunding := models.CrowdFunding{}
	filePath := ""
	if imageFile != nil && imageFileHeader != nil {
		uploadImageFolder := setting.CurrentConfig().UploadFolder + "/crowd-funding"
		fileName := imageFileHeader.Filename
		imageExt := strings.Split(fileName, ".")[1]
		imageName := "cf-" + time.Now().Format("20060102150405")
		fileNameImage := imageName + "." + imageExt
		err := fileUploadService.UploadFormFile(*imageFile, uploadImageFolder, fileNameImage, imageFileHeader)
		if err != nil {
			log.Println(err)
			return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
		}
		filePath = uploadImageFolder + "/" + fileNameImage
	}

	crowdFunding.UserId = userId
	crowdFunding.Name = request.Name
	crowdFunding.Description = request.Description
	crowdFunding.ShortDescription = request.ShortDescription
	crowdFunding.Image = filePath
	crowdFunding.YoutubeUrl = request.YoutubeUrl
	crowdFunding.CrowdDate = request.CrowdDate
	crowdFunding.DeliverDate = request.DeliverDate
	crowdFunding.Price = request.Price
	crowdFunding.Goal = request.Goal
	crowdFunding.Status = utils.CROWD_STATUS_NEW

	crowdFunding, err := crowdFundingDao.Create(crowdFunding, nil)
	if err != nil {
		log.Println(err)
		return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}
	return crowdFunding, nil
}

func (crowdService CrowdService) UpdateCrowdFunding(userId int64, crowdFundingId int64, request request_obj.CrowdFundingRequest, imageFile *multipart.File, imageFileHeader *multipart.FileHeader) (models.CrowdFunding, *bean.AppError) {
	crowdFunding := crowdFundingDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 || crowdFunding.UserId != userId {
		return crowdFunding, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	filePath := ""
	if imageFile != nil && imageFileHeader != nil {
		uploadImageFolder := setting.CurrentConfig().UploadFolder + "/crowd-funding"
		fileName := imageFileHeader.Filename
		imageExt := strings.Split(fileName, ".")[1]
		imageName := "cf-" + time.Now().Format("20060102150405")
		fileNameImage := imageName + "." + imageExt
		err := fileUploadService.UploadFormFile(*imageFile, uploadImageFolder, fileNameImage, imageFileHeader)
		if err != nil {
			log.Println(err)
			return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
		}
		filePath = uploadImageFolder + "/" + fileNameImage
	}

	crowdFunding.Name = request.Name
	crowdFunding.Description = request.Description
	crowdFunding.ShortDescription = request.ShortDescription
	if filePath != "" {
		crowdFunding.Image = filePath
	}
	if crowdFunding.Status == utils.CROWD_STATUS_NEW {
		crowdFunding.CrowdDate = request.CrowdDate
		crowdFunding.YoutubeUrl = request.YoutubeUrl
		crowdFunding.DeliverDate = request.DeliverDate
		crowdFunding.Price = request.Price
		crowdFunding.Goal = request.Goal
	}

	crowdFunding, err := crowdFundingDao.Update(crowdFunding, nil)
	if err != nil {
		log.Println(err)
		return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}
	return crowdFunding, nil
}

func (crowdService CrowdService) GetCrowdFunding(userId int64, crowdFundingId int64) (models.CrowdFunding, *bean.AppError) {
	crowdFunding := crowdFundingDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFunding, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	return crowdFunding, nil
}

func (crowdService CrowdService) ShakedCrowdFunding(userId int64, crowdFundingId int64, quantity int, address string, hash string) (models.CrowdFundingShaked, *bean.AppError) {
	crowdFundingShaked := models.CrowdFundingShaked{}

	if quantity <= 0 {
		return crowdFundingShaked, &bean.AppError{errors.New("quantity is invalid"), "quantity is invalid", -1, "error_occurred"}
	}

	crowdFunding := crowdFundingDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFundingShaked, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}

	crowdFundingShaked.UserId = userId
	crowdFundingShaked.CrowdFundingId = crowdFundingId
	crowdFundingShaked.Quantity = quantity
	crowdFundingShaked.Amount = float64(crowdFundingShaked.Quantity) * crowdFunding.Price
	crowdFundingShaked.Status = utils.CROWD_ORDER_STATUS_NEW

	crowdFundingShaked, err := crowdFundingShakedDao.Create(crowdFundingShaked, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingShaked, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	_, appErr := crowdService.CreateTx(userId, address, hash, "crowd_shake", crowdFundingShaked.ID, nil)
	if appErr != nil {
		log.Println(appErr.OrgError)
		return crowdFundingShaked, appErr
	}

	return crowdFundingShaked, nil
}
