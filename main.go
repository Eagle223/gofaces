package main

import (
	"fmt"
	"gofaces/router"
	"gofaces/rtsp"
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
	fmt.Println(imgDir)
	ticker := time.NewTicker(time.Hour * 1)
	defer ticker.Stop()
	for true {
		go rtsp.VideoCaptureStart1("classify%d.jpg", imgDir, ch)
		if 0 == strings.Compare("success", <-ch) {
			pid = <-ch
			fmt.Println("ffmpeg pid :", pid)
			imgDir = (<-ticker.C).String()
			fmt.Println(imgDir)
			rtsp.VideoCaptureStop1(pid, ch)
			<-ch
		}
	}
}
