package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-crowdfunding/bean"
	"github.com/ninjadotorg/handshake-crowdfunding/request_obj"
	"github.com/ninjadotorg/handshake-crowdfunding/response_obj"
)

type PostApi struct {
}

func (postApi PostApi) Init(router *gin.Engine) *gin.RouterGroup {
	faq := router.Group("/post")
	{
		faq.GET("/:crowd_funding_id", func(context *gin.Context) {
			context.String(200, "Common API")
		})
		faq.POST("/:crowd_funding_id", func(context *gin.Context) {
			postApi.CreatePost(context)
		})
		faq.PUT("/:faq_id", func(context *gin.Context) {
			postApi.CreatePost(context)
		})
	}
	return faq
}

func (postApi PostApi) CreatePost(context *gin.Context) {
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

	request := new(request_obj.CrowdFundingPostRequest)
	err = context.Bind(&request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	faq, err := crowdService.CreatePost(userId.(int64), crowdFungingId, *request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Data = response_obj.MakeCrowdFundingPostResponse(faq)
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (postApi PostApi) UpdatedUpdated(context *gin.Context) {
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
	crowdFungingUpdatedId, err := strconv.ParseInt(context.Param("faq_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if crowdFungingUpdatedId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	request := new(request_obj.CrowdFundingPostRequest)
	err = context.Bind(&request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	faq, err := crowdService.UpdatePost(userId.(int64), crowdFungingUpdatedId, *request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Data = response_obj.MakeCrowdFundingPostResponse(faq)
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}
