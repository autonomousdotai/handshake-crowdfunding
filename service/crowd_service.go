package service

import (
	"github.com/ninjadotorg/handshake-crowdfunding/models"
	"github.com/ninjadotorg/handshake-crowdfunding/bean"
	"errors"
	"mime/multipart"
	"strings"
	"time"
	"github.com/ninjadotorg/handshake-crowdfunding/configs"
	"log"
	"github.com/ninjadotorg/handshake-crowdfunding/request_obj"
	"github.com/ninjadotorg/handshake-crowdfunding/utils"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"github.com/rtt/Go-Solr"
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
			uploadImageFolder := "crowdfunding"
			fileName := imageFileHeader.Filename
			imageExt := strings.Split(fileName, ".")[1]
			fileNameImage := fmt.Sprintf("crowdfunding-%d-image-%s.%s", crowdFunding.ID, time.Now().Format("20060102150405"), imageExt)
			filePath = uploadImageFolder + "/" + fileNameImage
			err := fileUploadService.UploadFile(filePath, &imageFile)
			if err != nil {
				log.Println(err)
				//rollback
				tx.Rollback()
				return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
			}
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

func (crowdService CrowdService) UserShake(userId int64, crowdFundingId int64, quantity int, address string, hash string) (models.CrowdFundingShake, *bean.AppError) {
	crowdFundingShake := models.CrowdFundingShake{}

	if quantity <= 0 {
		return crowdFundingShake, &bean.AppError{errors.New("quantity is invalid"), "quantity is invalid", -1, "error_occurred"}
	}

	crowdFunding := crowdFundingDao.GetFullById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFundingShake, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}

	crowdFundingShake.UserId = userId
	crowdFundingShake.CrowdFundingId = crowdFundingId
	crowdFundingShake.Quantity = quantity
	crowdFundingShake.Amount = float64(crowdFundingShake.Quantity) * crowdFunding.Price
	crowdFundingShake.Status = utils.CROWD_ORDER_STATUS_NEW

	crowdFundingShake, err := crowdFundingShakeDao.Create(crowdFundingShake, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingShake, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	_, appErr := crowdService.CreateTx(userId, address, hash, "crowd_shake", crowdFundingShake.ID, nil)
	if appErr != nil {
		log.Println(appErr.OrgError)
		return crowdFundingShake, appErr
	}

	return crowdFundingShake, nil
}

func (crowdService CrowdService) UnshakeCrowdFunding(userId int64, crowdFundingId int64, address string, hash string) (*bean.AppError) {
	crowdFunding := crowdFundingDao.GetFullByUser(userId, crowdFundingId)
	if crowdFunding.ID <= 0 {
		return &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	crowdFundingShaked := crowdFunding.CrowdFundingShake

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
	crowdFundingShaked, err := crowdFundingShakeDao.Update(crowdFundingShaked, tx)
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
	crowdFundingShaked := crowdFunding.CrowdFundingShake

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
	crowdFundingShaked, err := crowdFundingShakeDao.Update(crowdFundingShaked, tx)
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
	crowdFundingShaked := crowdFunding.CrowdFundingShake

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
	crowdFundingShaked, err := crowdFundingShakeDao.Update(crowdFundingShaked, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService CrowdService) MakeObjectToIndex(crowdFundingId int64) (error) {
	crowdFunding := crowdFundingDao.GetFullById(crowdFundingId)

	crowdFundingImages := crowdFundingImageDao.GetByCrowdId(crowdFunding.ID)
	imageUrls := []string{}
	for _, crowdFundingImage := range crowdFundingImages {
		imageUrls = append(imageUrls, crowdFundingImage.Image)
	}

	document := map[string]interface{}{
		"add": [] interface{}{
			map[string]interface{}{
				"id":                fmt.Sprintf("crowd_%d", crowdFunding.ID),
				"hid_s":             "",
				"type_i":            1,
				"state_i":           0,
				"init_user_id_i":    crowdFunding.UserId,
				"shake_user_ids_is": []int64{},
				"text_search_ss":    []string{crowdFunding.Name, crowdFunding.Description, crowdFunding.ShortDescription},
				"shake_count_i":     crowdFunding.ShakeNum,
				"view_count_i":      0,
				"comment_count_i":   0,
				"is_private_i":      1,
				"init_at_i":         crowdFunding.DateCreated.Unix(),
				"last_update_at_i":  crowdFunding.DateModified.Unix(),
				//custom fileds
				"name_s":              crowdFunding.Name,
				"short_description_s": crowdFunding.ShortDescription,
				"goal_f":              crowdFunding.Goal,
				"balance_f":           crowdFunding.Balance,
				"crowd_date_i":        crowdFunding.CrowdDate.Unix(),
				"deliver_date_i":      crowdFunding.DeliverDate.Unix(),
				"image_ss":            imageUrls,
			},
		},
	}

	jsonStr, err := json.Marshal(document)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", configs.SolrServiceUrl+"/handshake/update", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	bodyBytes, err := netUtil.CurlRequest(req)
	if err != nil {
		return err
	}
	result := solr.UpdateResponse{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	if result.Success == false {
		return errors.New("update solr result false")
	}
	return nil
}

func (crowdService CrowdService) GetUser(userId int64) (models.User, error) {
	result := models.JsonUserResponse{}
	url := fmt.Sprintf("%s/%d", configs.DispatcherServiceUrl+"/system/user", userId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return result.Data, err
	}
	req.Header.Set("Content-Type", "application/json")
	bodyBytes, err := netUtil.CurlRequest(req)
	if err != nil {
		log.Println(err)
		return result.Data, err
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Println(err)
		return result.Data, err
	}
	return result.Data, err
}

func (crowdService CrowdService) CreateFaq(userId int64, crowdFundingId int64, crowdFundingFaqRequest request_obj.CrowdFundingFaqRequest) (models.CrowdFundingFaq, *bean.AppError) {
	crowdFundingFaq := models.CrowdFundingFaq{}

	crowdFundingFaq.UserId = userId
	crowdFundingFaq.CrowdFundingId = crowdFundingId
	crowdFundingFaq.Question = crowdFundingFaqRequest.Question
	crowdFundingFaq.Answer = crowdFundingFaqRequest.Answer
	crowdFundingFaq.Status = 1

	crowdFundingFaq, err := crowdFundingFaqDao.Create(crowdFundingFaq, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingFaq, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	return crowdFundingFaq, nil
}

func (crowdService CrowdService) GetFaqsByCrowdId(crowdFundingId int64, pagination *bean.Pagination) (*bean.Pagination, error) {
	pagination, err := crowdFundingFaqDao.GetAllBy(0, crowdFundingId, pagination)
	faqs := pagination.Items.([]models.CrowdFundingFaq)
	items := []models.CrowdFundingFaq{}
	for _, faq := range faqs {
		user, _ := crowdService.GetUser(faq.UserId)
		faq.User = user
		items = append(items, faq)
	}
	pagination.Items = items
	return pagination, err
}
