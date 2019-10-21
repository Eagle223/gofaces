package main

import (
	"gofaces/imageprocess"
	"gofaces/rtsp"
	"log"
)

func main() {
	ch := make(chan string)
	go rtsp.VideoCaptureHandler1()

	if 3 == imageprocess.BuildClassifier() {
		log.Println("加载Data 错误！")
		return
	}
	go imageprocess.FaceDetections(ch)
	for true {
		log.Println("main:" + <-ch)
	}
}
