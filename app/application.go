package app

import (
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default() //private variable available only in application
)

func StartApplication() {
	mapUrls()
	//gin.SetMode(gin.ReleaseMode)
	router.Run(":8080")
}
