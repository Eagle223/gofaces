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
	api.GET("/startVideoCapture", rtsp.VideoCaptureStart)
	api.GET("/stopVideoCapture", rtsp.VideoCaptureStop)
	api.GET("/buildFaceModle", imageprocess.BuildFaceModle)
	api.GET("/classifyFace", imageprocess.ClassifyFace)
	Router.Run()

}
