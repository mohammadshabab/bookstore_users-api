package app

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammadshabab/bookstore_users-api/logger"
)

var (
	router = gin.Default() //private variable available only in application
)

func StartApplication() {
	mapUrls()
	//gin.SetMode(gin.ReleaseMode)
	logger.Info("about to start the application ...")
	router.Run(":8080")
}
