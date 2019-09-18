package main

import (
	"gofaces/router"
	"gofaces/rtsp"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	go videoCaptureHandler()
	router.Init()
}

func videoCaptureHandler() {
	var pid string
	ch := make(chan string)
	imgDir := strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Year())
	log.Println(imgDir)
	go rtsp.VideoCaptureStart1("classify%d.jpg", imgDir, ch)
	if 0 == strings.Compare("success", <-ch) {
		ticker1 := time.NewTicker(time.Minute * time.Duration(60-time.Now().Minute()))
		pid = <-ch
		log.Println("ticker time out", <-ticker1.C)
		go rtsp.VideoCaptureStop1(pid, ch)
		<-ch
		ticker1.Stop()
	} else {
		log.Fatalf("ffmpeg 出错，请检查系统")
	}
	ticker := time.NewTicker(time.Hour)
	for true {
		imgDir = strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Year())
		log.Println(imgDir)
		go rtsp.VideoCaptureStart1("classify%d.jpg", imgDir, ch)
		if 0 == strings.Compare("success", <-ch) {
			pid = <-ch
			log.Println("ffmpeg pid :", pid)
			log.Println("ticker time out", <-ticker.C)
			go rtsp.VideoCaptureStop1(pid, ch)
			<-ch
			rtsp.CleanOldImages()
		} else {
			log.Fatalf("ffmpeg 出错，请检查系统")
		}

	}
}
