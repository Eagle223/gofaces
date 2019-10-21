package router

import (
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.New()
	Router.Run()
}
