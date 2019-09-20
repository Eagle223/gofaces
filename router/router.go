package router

import (
	"github.com/gin-gonic/gin"
	"gofaces/imageprocess"
	"gofaces/rtsp"
)

var Router *gin.Engine

func Init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.New()
	api := Router.Group("api/v1")
	api.GET("/buildFaceModle", imageprocess.BuildFaceModle)
	api.GET("/classifyFace", imageprocess.ClassifyFace)
	api.GET("/getLatestImage", rtsp.GetLatestImage1)
	Router.Run()
}
