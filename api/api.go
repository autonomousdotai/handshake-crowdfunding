package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-crowdfunding/response_obj"
	"github.com/ninjadotorg/handshake-crowdfunding/request_obj"
	"net/http"
	"strconv"
	"github.com/ninjadotorg/handshake-crowdfunding/bean"
	"log"
	"encoding/json"
)

type Api struct {
}

func (api Api) Init(router *gin.Engine) *gin.Engine {
	router.POST("/", func(context *gin.Context) {
		api.CreateCrowdFunding(context)
	})
	router.PUT("/", func(context *gin.Context) {
		api.UpdateCrowdFunding(context)
	})
	router.GET("/:crowd_funding_id", func(context *gin.Context) {
		api.GetCrowdFunding(context)
	})
	router.POST("/shake/:crowd_funding_id", func(context *gin.Context) {
		api.UserShake(context)
	})
	return router
}

func (self Api) CreateCrowdFunding(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}

	requestJson := context.Request.PostFormValue("request")
	request := new(request_obj.CrowdFundingRequest)
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	crowdFunging, appErr := crowdService.CreateCrowdFunding(userId.(int64), *request, context)
	if appErr != nil {
		log.Print(appErr.OrgError)
		result.SetStatus(bean.UnexpectedError)
		result.Error = appErr.OrgError.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	data := response_obj.MakeCrowdFundingResponse(crowdFunging)

	result.Data = data
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self Api) UpdateCrowdFunding(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	crowdFungingId, err := strconv.ParseInt(context.Param("crowd_funding_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if crowdFungingId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	requestJson := context.Request.PostFormValue("request")
	request := new(request_obj.CrowdFundingRequest)
	err = json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	imageFile, imageFileHeader, err := context.Request.FormFile("image")
	crowdFunging, appErr := crowdService.UpdateCrowdFunding(userId.(int64), crowdFungingId, *request, &imageFile, imageFileHeader)
	if appErr != nil {
		log.Print(appErr.OrgError)
		result.SetStatus(bean.UnexpectedError)
		result.Error = appErr.OrgError.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	data := response_obj.MakeCrowdFundingResponse(crowdFunging)

	result.Data = data
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self Api) GetCrowdFunding(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	crowdFungingId, err := strconv.ParseInt(context.Param("crowd_funding_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if crowdFungingId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	crowdFunging, appErr := crowdService.GetCrowdFunding(0, crowdFungingId)
	if appErr != nil {
		log.Print(appErr.OrgError)
		result.SetStatus(bean.UnexpectedError)
		result.Error = appErr.OrgError.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	data := response_obj.MakeCrowdFundingResponse(crowdFunging)

	result.Data = data
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self Api) UserShake(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	crowdFungingId, err := strconv.ParseInt(context.Param("crowd_funding_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if crowdFungingId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	quantity, err := strconv.Atoi(context.Param("quantity"))
	address := context.Query("address")
	hash := context.Query("hash")

	crowdFungingShaked, appErr := crowdService.UserShake(userId.(int64), crowdFungingId, quantity, address, hash)
	if appErr != nil {
		log.Print(appErr.OrgError)
		result.SetStatus(bean.UnexpectedError)
		result.Error = appErr.OrgError.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	_ = crowdFungingShaked

	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}
