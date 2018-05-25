package service

import (
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/models"
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/bean"
	"errors"
	"mime/multipart"
	"strings"
	"time"
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/setting"
	"log"
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/request_obj"
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/utils"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
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

func (crowdService CrowdService) CreateCrowdFunding(userId int64, request request_obj.CrowdFundingRequest, context *gin.Context) (models.CrowdFunding, *bean.AppError) {
	crowdFunding := models.CrowdFunding{}

	tx := models.Database().Begin()

	crowdFunding.UserId = userId
	crowdFunding.Name = request.Name
	crowdFunding.Description = request.Description
	crowdFunding.ShortDescription = request.ShortDescription
	crowdFunding.CrowdDate = request.CrowdDate
	crowdFunding.DeliverDate = request.DeliverDate
	crowdFunding.Price = request.Price
	crowdFunding.Goal = request.Goal
	crowdFunding.Status = utils.CROWD_STATUS_NEW

	crowdFunding, err := crowdFundingDao.Create(crowdFunding, tx)
	if err != nil {
		log.Println(err)
		//rollback
		tx.Rollback()
		return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	imageLength, err := strconv.Atoi(context.Request.PostFormValue("image_length"))
	for i := 0; i < imageLength; i++ {
		imageFile, imageFileHeader, err := context.Request.FormFile(fmt.Sprintf("image_%d", i))
		if err != nil {
			log.Println(err)
			//rollback
			tx.Rollback()
			return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
		}
		filePath := ""
		if imageFile != nil && imageFileHeader != nil {
			uploadImageFolder := setting.CurrentConfig().UploadFolder + "/crowd-funding/image"
			fileName := imageFileHeader.Filename
			imageExt := strings.Split(fileName, ".")[1]
			fileNameImage := fmt.Sprintf("crowd-funding-%d-image-%s.%s", crowdFunding.ID, time.Now().Format("20060102150405"), imageExt)
			err := fileUploadService.UploadFormFile(imageFile, uploadImageFolder, fileNameImage, imageFileHeader)
			if err != nil {
				log.Println(err)
				//rollback
				tx.Rollback()
				return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
			}
			filePath = uploadImageFolder + "/" + fileNameImage
		}
		crowdFundingImage := models.CrowdFundingImage{}

		crowdFundingImage.CrowdFundingId = crowdFunding.ID
		crowdFundingImage.Image = filePath

		crowdFundingImage, err = crowdFundingImageDao.Create(crowdFundingImage, tx)
		if err != nil {
			log.Println(err)
			//rollback
			tx.Rollback()
			return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
		}
	}

	tx.Commit()

	crowdFunding = crowdFundingDao.GetFullById(crowdFunding.ID)

	return crowdFunding, nil
}

func (crowdService CrowdService) UpdateCrowdFunding(userId int64, crowdFundingId int64, request request_obj.CrowdFundingRequest, imageFile *multipart.File, imageFileHeader *multipart.FileHeader) (models.CrowdFunding, *bean.AppError) {
	crowdFunding := crowdFundingDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 || crowdFunding.UserId != userId {
		return crowdFunding, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}

	crowdFunding.Name = request.Name
	crowdFunding.Description = request.Description
	crowdFunding.ShortDescription = request.ShortDescription

	if crowdFunding.Status == utils.CROWD_STATUS_NEW {
		crowdFunding.CrowdDate = request.CrowdDate
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

	crowdFunding := crowdFundingDao.GetFullById(crowdFundingId)
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

func (crowdService CrowdService) UnshakeCrowdFunding(userId int64, crowdFundingId int64, address string, hash string) (*bean.AppError) {
	crowdFunding := crowdFundingDao.GetFullByUser(userId, crowdFundingId)
	if crowdFunding.ID <= 0 {
		return &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	crowdFundingShaked := crowdFunding.CrowdFundingShaked

	if crowdFundingShaked.ID <= 0 || crowdFundingShaked.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}

	tx := models.Database().Begin()

	_, appErr := crowdService.CreateTx(userId, address, hash, "crowd_unshake", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}

	crowdFundingShaked.Status = utils.CROWD_ORDER_STATUS_UNSHAKED_PROCESS
	crowdFundingShaked, err := crowdFundingShakedDao.Update(crowdFundingShaked, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService CrowdService) CancelCrowdFunding(userId int64, crowdFundingId int64, address string, hash string) (*bean.AppError) {
	crowdFunding := crowdFundingDao.GetFullByUser(userId, crowdFundingId)
	if crowdFunding.ID <= 0 {
		return &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	crowdFundingShaked := crowdFunding.CrowdFundingShaked

	if crowdFundingShaked.ID <= 0 || crowdFundingShaked.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}

	tx := models.Database().Begin()

	_, appErr := crowdService.CreateTx(userId, address, hash, "crowd_cancel", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}

	crowdFundingShaked.Status = utils.CROWD_ORDER_STATUS_CANCELED_PROCESS
	crowdFundingShaked, err := crowdFundingShakedDao.Update(crowdFundingShaked, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService CrowdService) RefundCrowdFunding(userId int64, crowdFundingId int64, address string, hash string) (*bean.AppError) {
	crowdFunding := crowdFundingDao.GetFullByUser(userId, crowdFundingId)
	if crowdFunding.ID <= 0 {
		return &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	crowdFundingShaked := crowdFunding.CrowdFundingShaked

	if crowdFundingShaked.ID <= 0 || crowdFundingShaked.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}

	tx := models.Database().Begin()

	_, appErr := crowdService.CreateTx(userId, address, hash, "crowd_refund", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}

	crowdFundingShaked.Status = utils.CROWD_ORDER_STATUS_REFUNDED_PROCESS
	crowdFundingShaked, err := crowdFundingShakedDao.Update(crowdFundingShaked, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}
