package api

import (
	"github.com/gin-gonic/gin"
	"../response_obj"
	"../request_obj"
	"net/http"
	"strconv"
	"../bean"
	"log"
	"encoding/json"
)

type Api struct {
}

func (api Api) Init(router *gin.Engine) *gin.RouterGroup {
	apiGroupApi := router.Group("/api")
	{
		apiGroupApi.GET("/", func(context *gin.Context) {
			context.String(200, "Crowdsale API")
		})
		apiGroupApi.POST("/crowd-funding", func(context *gin.Context) {
			api.CreateCrowdFunding(context)
		})
		apiGroupApi.PUT("/crowd-funding", func(context *gin.Context) {
			api.UpdateCrowdFunding(context)
		})
		apiGroupApi.GET("/crowd-funding/:crowd_funding_id", func(context *gin.Context) {
			api.GetCrowdFunding(context)
		})
		apiGroupApi.POST("/crowd-funding/shaked/:crowd_funding_id", func(context *gin.Context) {
			api.ShakedCrowdFunding(context)
		})
	}
	return apiGroupApi
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
	imageFile, imageFileHeader, err := context.Request.FormFile("image")
	crowdFunging, appErr := crowdService.CreateCrowdFunding(userId.(int64), *request, &imageFile, imageFileHeader)
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

func (self Api) ShakedCrowdFunding(context *gin.Context) {
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

	crowdFungingShaked, appErr := crowdService.ShakedCrowdFunding(userId.(int64), crowdFungingId, quantity, address, hash)
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
