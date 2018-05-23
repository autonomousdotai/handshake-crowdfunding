package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"../api"
	"../setting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
)

func main() {

	configuration := setting.CurrentConfig()
	// Logger
	logFile, err := os.OpenFile("logs/autonomous_service.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter) // You may need this
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// end Logger
	// Setting router
	router := gin.Default()
	router.Use(Logger())
	router.Use(sentry.Recovery(raven.DefaultClient, false))
	router.Use(CORSMiddleware())
	router.Use(AuthorizeMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// Router Index
	index := router.Group("/")
	{
		index.GET("/", func(context *gin.Context) {
			context.String(http.StatusOK, "API")
		})
	}
	api := api.Api{}
	api.Init(router)
	router.Run(fmt.Sprintf(":%d", configuration.ServicePort))
}

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		context.Next()
		status := context.Writer.Status()
		latency := time.Since(t)
		log.Print("Request: " + context.Request.URL.String() + " | " + context.Request.Method + " - Status: " + strconv.Itoa(status) + " - " +
			latency.String())
	}
}

func CORSMiddleware() gin.HandlerFunc {
	// Gin Cors setting
	return cors.New(cors.Config{
		AllowOrigins:     []string{setting.CurrentConfig().OriginDomainLocal},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"},
		AllowHeaders:     []string{"Content-Type", "Origin", "Device-Type", "Device-Id", "*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, _ := strconv.ParseInt(context.GetHeader("User-Id"), 10, 64)
		context.Set("UserId", userId)
		context.Next()
	}
}
