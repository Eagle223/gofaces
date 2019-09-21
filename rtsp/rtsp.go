package rtsp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"time"
)

const RtspUrl = "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream"
const ImgRootUrl = "/home/gofaces/rtmp/"

var classifyPID int

func VideoCaptureStart1(imgName string, imgDir string, ch chan<- string) {
	commandFmt := "ffmpeg  -y -i \"%v\" -r 1 %v"
	img := ImgRootUrl + imgDir
	log.Println("mkdir -p ", img)
	mkdir := exec.Command("mkdir", "-p", img)
	err := mkdir.Run()
	if err != nil {
		log.Println("mkdir erro:%v", err)
		ch <- fmt.Sprintf("mkdir erro:%v", err)
		return
	}
	img = img + "/" + imgName
	command := fmt.Sprintf(commandFmt, RtspUrl, img)
	log.Println(command)
	cmd := exec.Command("sh", "-c", command)
	err = cmd.Start()
	if err != nil {
		log.Println("ffmpeg start error:%v", err)
		ch <- fmt.Sprint("ffmpeg start err:%v", err)
		return
	}
	ch <- "success"
	ch <- strconv.Itoa(cmd.Process.Pid + 1)
	cmd.Wait()
	mkdir.Wait()
}

func VideoCaptureStop1(pid string, ch chan<- string) {
	cmd := exec.Command("sh", "-c", "kill "+pid)
	err := cmd.Run()
	if err != nil {
		log.Println("ffmpeg close error:%v", err)
		ch <- fmt.Sprintf("ffmpeg close error:%v", err)
	} else {
		ch <- "success"
		log.Println("ffmpe kill finished")
	}
	cmd.Wait()
}

func GetLatestImage() string {
	timeNow := time.Now()
	imgDirNow := strconv.Itoa(timeNow.Hour()) + "-" + strconv.Itoa(timeNow.Day()) + "-" + timeNow.Month().String() + "-" + strconv.Itoa(timeNow.Year())
	files, _ := ioutil.ReadDir(ImgRootUrl + imgDirNow + "/")
	rec := "error"
	if len(files) > 0 {
		log.Println("maxNum:", len(files))
		rec = ImgRootUrl + imgDirNow + "/" + "classify" + strconv.Itoa(len(files)) + ".jpg"
	} else {
		timeLast := timeNow.Add(-time.Hour)
		imgDirLast := strconv.Itoa(timeLast.Hour()) + "-" + strconv.Itoa(timeLast.Day()) + "-" + timeLast.Month().String() + "-" + strconv.Itoa(timeLast.Year())
		files, _ = ioutil.ReadDir(ImgRootUrl + imgDirLast + "/")
		if len(files) > 0 {
			log.Println("maxNum:", len(files))
			rec = ImgRootUrl + imgDirLast + "/" + "classify" + strconv.Itoa(len(files)) + ".jpg"
		}
	}
	return rec
}

func CleanOldImages() {
	timeOld := time.Now().Add(-time.Hour * 3)
	imgDirOld := strconv.Itoa(timeOld.Hour()) + "-" + strconv.Itoa(timeOld.Day()) + "-" + timeOld.Month().String() + "-" + strconv.Itoa(timeOld.Year())
	cmd := exec.Command("rm", "-rf", ImgRootUrl+imgDirOld)
	err := cmd.Run()
	if err != nil {
		log.Println("clean old Images error:", err)
	}
	log.Println("clean:", ImgRootUrl+imgDirOld)
	cmd.Wait()
}

func GetLatestImage1(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": GetLatestImage(),
	})
}
