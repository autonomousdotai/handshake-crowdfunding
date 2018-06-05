package api

import (
	"github.com/gin-gonic/gin"
)

type FaqApi struct {
}

func (faqApi FaqApi) Init(router *gin.Engine) *gin.RouterGroup {
	faq := router.Group("/faq")
	{
		faq.GET("/:faq_id", func(context *gin.Context) {
			context.String(200, "Common API")
		})
	}
	return faq
}
