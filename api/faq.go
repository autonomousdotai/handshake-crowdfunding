package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-crowdfunding/response_obj"
	"github.com/ninjadotorg/handshake-crowdfunding/bean"
	"net/http"
	"strconv"
	"github.com/ninjadotorg/handshake-crowdfunding/request_obj"
	"log"
)

type FaqApi struct {
}

func (faqApi FaqApi) Init(router *gin.Engine) *gin.RouterGroup {
	faq := router.Group("/faq")
	{
		faq.GET("/:crowd_funding_id", func(context *gin.Context) {
			context.String(200, "Common API")
		})
		faq.POST("/:crowd_funding_id", func(context *gin.Context) {
			faqApi.CreateFaq(context)
		})
		faq.PUT("/:faq_id", func(context *gin.Context) {
			faqApi.CreateFaq(context)
		})
	}
	return faq
}

func (faqApi FaqApi) CreateFaq(context *gin.Context) {
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

	request := new(request_obj.CrowdFundingFaqRequest)
	err = context.Bind(&request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	faq, err := crowdService.CreateFaq(userId.(int64), crowdFungingId, *request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Data = response_obj.MakeCrowdFundingFaqResponse(faq)
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (faqApi FaqApi) UpdateFaq(context *gin.Context) {
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
	crowdFungingFaqId, err := strconv.ParseInt(context.Param("faq_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if crowdFungingFaqId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	request := new(request_obj.CrowdFundingFaqRequest)
	err = context.Bind(&request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	faq, err := crowdService.UpdateFaq(userId.(int64), crowdFungingFaqId, *request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Data = response_obj.MakeCrowdFundingFaqResponse(faq)
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}
