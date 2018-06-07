package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-crowdfunding/bean"
	"github.com/ninjadotorg/handshake-crowdfunding/configs"
	"github.com/ninjadotorg/handshake-crowdfunding/models"
	"github.com/ninjadotorg/handshake-crowdfunding/request_obj"
	"github.com/ninjadotorg/handshake-crowdfunding/utils"
	"github.com/rtt/Go-Solr"
)

type CrowdService struct {
}

func (crowdService CrowdService) CreateCrowdFunding(userId int64, request request_obj.CrowdFundingRequest, context *gin.Context) (models.CrowdFunding, error) {
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
		return crowdFunding, err
	}

	imageLength, err := strconv.Atoi(context.Request.PostFormValue("image_length"))
	for i := 0; i < imageLength; i++ {
		imageFile, imageFileHeader, err := context.Request.FormFile(fmt.Sprintf("image_%d", i))
		if err != nil {
			log.Println(err)
			//rollback
			tx.Rollback()
			return crowdFunding, err
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
				return crowdFunding, err
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
			return crowdFunding, err
		}
	}

	tx.Commit()

	crowdFunding = crowdFundingDao.GetFullById(crowdFunding.ID)

	return crowdFunding, nil
}

func (crowdService CrowdService) UpdateCrowdFunding(userId int64, crowdFundingId int64, request request_obj.CrowdFundingRequest, imageFile *multipart.File, imageFileHeader *multipart.FileHeader) (models.CrowdFunding, error) {
	crowdFunding := crowdFundingDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 || crowdFunding.UserId != userId {
		return crowdFunding, errors.New("crowdFundingId is invalid")
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
		return crowdFunding, err
	}
	return crowdFunding, nil
}

func (crowdService CrowdService) GetCrowdFunding(userId int64, crowdFundingId int64) (models.CrowdFunding, error) {
	crowdFunding := crowdFundingDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFunding, errors.New("crowdFundingId is invalid")
	}
	return crowdFunding, nil
}

func (crowdService CrowdService) UserShake(userId int64, crowdFundingId int64, quantity int, address string, hash string) (models.CrowdFundingShake, error) {
	crowdFundingShake := models.CrowdFundingShake{}

	if quantity <= 0 {
		return crowdFundingShake, errors.New("quantity is invalid")
	}

	crowdFunding := crowdFundingDao.GetFullById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFundingShake, errors.New("crowdFundingId is invalid")
	}

	crowdFundingShake.UserId = userId
	crowdFundingShake.CrowdFundingId = crowdFundingId
	crowdFundingShake.Quantity = quantity
	crowdFundingShake.Amount = float64(crowdFundingShake.Quantity) * crowdFunding.Price
	crowdFundingShake.Status = utils.CROWD_ORDER_STATUS_NEW

	crowdFundingShake, err := crowdFundingShakeDao.Create(crowdFundingShake, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingShake, err
	}

	return crowdFundingShake, nil
}

func (crowdService CrowdService) UnshakeCrowdFunding(userId int64, crowdFundingId int64) error {
	crowdFunding := crowdFundingDao.GetFullByUser(userId, crowdFundingId)
	if crowdFunding.ID <= 0 {
		return errors.New("crowdFundingId is invalid")
	}
	crowdFundingShaked := crowdFunding.CrowdFundingShake

	if crowdFundingShaked.ID <= 0 || crowdFundingShaked.Status <= 0 {
		return errors.New("crowdFundingId is invalid")
	}

	crowdFundingShaked.Status = utils.CROWD_ORDER_STATUS_UNSHAKED_PROCESS
	crowdFundingShaked, err := crowdFundingShakeDao.Update(crowdFundingShaked, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (crowdService CrowdService) CancelCrowdFunding(userId int64, crowdFundingId int64) error {
	crowdFunding := crowdFundingDao.GetFullByUser(userId, crowdFundingId)
	if crowdFunding.ID <= 0 {
		return &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	crowdFundingShaked := crowdFunding.CrowdFundingShake

	if crowdFundingShaked.ID <= 0 || crowdFundingShaked.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}

	crowdFundingShaked.Status = utils.CROWD_ORDER_STATUS_CANCELED_PROCESS
	crowdFundingShaked, err := crowdFundingShakeDao.Update(crowdFundingShaked, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (crowdService CrowdService) RefundCrowdFunding(userId int64, crowdFundingId int64) error {
	crowdFunding := crowdFundingDao.GetFullByUser(userId, crowdFundingId)
	if crowdFunding.ID <= 0 {
		return &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	crowdFundingShaked := crowdFunding.CrowdFundingShake

	if crowdFundingShaked.ID <= 0 || crowdFundingShaked.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}

	crowdFundingShaked.Status = utils.CROWD_ORDER_STATUS_REFUNDED_PROCESS
	crowdFundingShaked, err := crowdFundingShakeDao.Update(crowdFundingShaked, nil)
	if err != nil {
		log.Println(err)

		return err
	}

	return nil
}

func (crowdService CrowdService) IndexSolr(crowdFundingId int64) error {
	crowdFunding := crowdFundingDao.GetFullById(crowdFundingId)

	crowdFundingImages := crowdFundingImageDao.GetByCrowdId(crowdFunding.ID)
	imageUrls := []string{}
	for _, crowdFundingImage := range crowdFundingImages {
		imageUrls = append(imageUrls, crowdFundingImage.Image)
	}

	document := map[string]interface{}{
		"add": []interface{}{
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
	url := fmt.Sprintf("%s/%s", configs.AppConf.SolrServiceUrl, "handshake/update")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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
	url := fmt.Sprintf("%s/%s/%d", configs.AppConf.DispatcherServiceUrl, "system/user", userId)
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

func (crowdService CrowdService) CreateFaq(userId int64, crowdFundingId int64, crowdFundingFaqRequest request_obj.CrowdFundingFaqRequest) (models.CrowdFundingFaq, error) {
	crowdFundingFaq := models.CrowdFundingFaq{}

	crowdFundingFaq.UserId = userId
	crowdFundingFaq.CrowdFundingId = crowdFundingId
	crowdFundingFaq.Question = crowdFundingFaqRequest.Question
	crowdFundingFaq.Answer = crowdFundingFaqRequest.Answer
	crowdFundingFaq.Status = 1

	crowdFundingFaq, err := crowdFundingFaqDao.Create(crowdFundingFaq, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingFaq, err
	}

	return crowdFundingFaq, nil
}

func (crowdService CrowdService) UpdateFaq(userId int64, faqId int64, crowdFundingFaqRequest request_obj.CrowdFundingFaqRequest) (models.CrowdFundingFaq, error) {
	crowdFundingFaq := crowdFundingFaqDao.GetById(faqId)

	if crowdFundingFaq.ID <= 0 || crowdFundingFaq.UserId != userId {
		return crowdFundingFaq, errors.New("faq_id is invalid")
	}

	crowdFundingFaq.Question = crowdFundingFaqRequest.Question
	crowdFundingFaq.Answer = crowdFundingFaqRequest.Answer

	crowdFundingFaq, err := crowdFundingFaqDao.Update(crowdFundingFaq, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingFaq, err
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

func (crowdService CrowdService) ProcessEventInit(hid int64, crowdFungdingId int64) error {
	crowdFunding := crowdFundingDao.GetById(crowdFungdingId)
	if crowdFunding.ID <= 0 {
		return errors.New("crowdFunding is invalid")
	}

	crowdFunding.Hid = hid
	crowdFunding.Status = utils.CROWD_STATUS_APPROVED

	crowdFunding, err := crowdFundingDao.Update(crowdFunding, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (crowdService CrowdService) ProcessEventShake(hid int64, state int, balance float64, crowdFundingShakeId int64, fromAddress string) error {
	crowdFundingShake := crowdFundingShakeDao.GetById(crowdFundingShakeId)
	if crowdFundingShake.ID < 0 {
		return errors.New("crowdFundingShake is invalid")
	}
	tx := models.Database().Begin()
	//for check shaked before
	crowdFundingShakes := crowdFundingShakeDao.GetAllByBackerStatus(crowdFundingShake.CrowdFundingId, crowdFundingShake.UserId, utils.CROWD_ORDER_STATUS_SHAKED)

	crowdFundingShake.Address = fromAddress
	crowdFundingShake.Status = utils.CROWD_ORDER_STATUS_SHAKED

	crowdFundingShake, err := crowdFundingShakeDao.Update(crowdFundingShake, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	crowdFunding := crowdFundingDao.GetByHId(hid)
	if crowdFunding.ID <= 0 {
		return errors.New("crowdFunding is invalid")
	}

	crowdFunding.Balance = balance / math.Pow(10, 18)
	if len(crowdFundingShakes) == 0 {
		crowdFunding.ShakeNum += 1
	}

	if state == 1 && crowdFunding.Status == utils.CROWD_STATUS_FAILED && crowdFunding.CrowdDate.Before(time.Now()) {
		crowdFunding.Status = utils.CROWD_STATUS_FUNDED
	}

	crowdFunding, err = crowdFundingDao.Update(crowdFunding, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService CrowdService) ProcessEventUnShake(hid int64, state int, balance float64, userId int64) error {
	tx := models.Database().Begin()

	crowdFunding := crowdFundingDao.GetByHId(hid)
	if crowdFunding.ID <= 0 {
		return errors.New("crowdFunding is invalid")
	}

	crowdFundingShakes := crowdFundingShakeDao.GetAllByBackerStatus(crowdFunding.ID, userId, utils.CROWD_ORDER_STATUS_UNSHAKED_PROCESS)
	for _, crowdFundingShake := range crowdFundingShakes {
		crowdFundingShake.Status = utils.CROWD_ORDER_STATUS_UNSHAKED
		_, err := crowdFundingShakeDao.Update(crowdFundingShake, tx)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
	}

	crowdFunding.Balance = balance / math.Pow(10, 18)
	crowdFunding.ShakeNum -= 1

	crowdFunding, err := crowdFundingDao.Update(crowdFunding, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService CrowdService) ProcessEventCancel(hid int64, state int, userId int64) error {
	tx := models.Database().Begin()

	crowdFunding := crowdFundingDao.GetByHId(hid)
	if crowdFunding.ID <= 0 {
		return errors.New("crowdFunding is invalid")
	}

	crowdFundingShakes := crowdFundingShakeDao.GetAllByBackerStatus(crowdFunding.ID, userId, utils.CROWD_ORDER_STATUS_CANCELED_PROCESS)
	for _, crowdFundingShake := range crowdFundingShakes {
		crowdFundingShake.Status = utils.CROWD_ORDER_STATUS_CANCELED
		_, err := crowdFundingShakeDao.Update(crowdFundingShake, tx)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
	}

	if state == 2 {
		crowdFunding.Status = utils.CROWD_STATUS_CANCELED
		_, err := crowdFundingDao.Update(crowdFunding, tx)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (crowdService CrowdService) ProcessEventRefund(hid int64, state int, userId int64) error {
	tx := models.Database().Begin()

	crowdFunding := crowdFundingDao.GetByHId(hid)
	if crowdFunding.ID <= 0 {
		return errors.New("crowdFunding is invalid")
	}

	crowdFundingShakes := crowdFundingShakeDao.GetAllByBackerStatus(crowdFunding.ID, userId, utils.CROWD_ORDER_STATUS_REFUNDED_PROCESS)
	for _, crowdFundingShake := range crowdFundingShakes {
		crowdFundingShake.Status = utils.CROWD_ORDER_STATUS_REFUNDED
		_, err := crowdFundingShakeDao.Update(crowdFundingShake, tx)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (crowdService CrowdService) ProcessEventStop(hid int64, state int, crowdFungdingId int64) error {
	crowdFunding := crowdFundingDao.GetById(crowdFungdingId)
	if crowdFunding.ID <= 0 {
		return errors.New("crowdFunding is invalid")
	}

	crowdFunding.Status = utils.CROWD_STATUS_CANCELED

	crowdFunding, err := crowdFundingDao.Update(crowdFunding, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (crowdService CrowdService) ProcessEventWithdraw(hid int64, amount float64, userId int64) error {
	//email withdraw refund amount successful
	return nil
}

func (crowdService CrowdService) CreatePost(userId int64, crowdFundingId int64, crowdFundingPostRequest request_obj.CrowdFundingPostRequest) (models.CrowdFundingPost, error) {
	crowdFundingPost := models.CrowdFundingPost{}

	crowdFundingPost.UserId = userId
	crowdFundingPost.CrowdFundingId = crowdFundingId
	crowdFundingPost.Title = crowdFundingPostRequest.Title
	crowdFundingPost.ShortDescription = crowdFundingPostRequest.ShortDescription
	crowdFundingPost.Description = crowdFundingPostRequest.Description
	crowdFundingPost.Status = 1

	crowdFundingPost, err := crowdFundingPostDao.Create(crowdFundingPost, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingPost, err
	}

	return crowdFundingPost, nil
}

func (crowdService CrowdService) UpdatePost(userId int64, updatedId int64, crowdFundingPostRequest request_obj.CrowdFundingPostRequest) (models.CrowdFundingPost, error) {
	crowdFundingPost := crowdFundingPostDao.GetById(updatedId)

	if crowdFundingPost.ID <= 0 || crowdFundingPost.UserId != userId {
		return crowdFundingPost, errors.New("post_id is invalid")
	}

	crowdFundingPost.Title = crowdFundingPostRequest.Title
	crowdFundingPost.ShortDescription = crowdFundingPostRequest.ShortDescription
	crowdFundingPost.Description = crowdFundingPostRequest.Description

	crowdFundingPost, err := crowdFundingPostDao.Update(crowdFundingPost, nil)
	if err != nil {
		log.Println(err)
		return crowdFundingPost, err
	}

	return crowdFundingPost, nil
}

func (crowdService CrowdService) GetPostsByCrowdId(crowdFundingId int64, pagination *bean.Pagination) (*bean.Pagination, error) {
	pagination, err := crowdFundingPostDao.GetAllBy(0, crowdFundingId, pagination)
	faqs := pagination.Items.([]models.CrowdFundingPost)
	items := []models.CrowdFundingPost{}
	for _, faq := range faqs {
		user, _ := crowdService.GetUser(faq.UserId)
		faq.User = user
		items = append(items, faq)
	}
	pagination.Items = items
	return pagination, err
}
