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

//ffmpeg -i "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream" -r 1  -y /home/eagle/rtmp/images/image%d.jpg 每秒截取一张图片
//ffmpeg -y -i "rtmp://58.200.131.2:1935/livetv/hunantv" -ss 00:00:01 -vframes 1 -f image2 /home/eagle/rtmp/images/image1.jpg

//const RtspUrl = "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream"
const RtspUrl = "rtsp://admin:123456@192.168.0.221:554/h264/ch1/main/av_stream"
const ImgRootUrl = "/home/eagle/rtmp/"

var classifyPID int

/*
* ffmpeg 参数设置
* -i 输入流或者文件
* -r 设置帧数Hz为单位
* -s 设置帧大小（WxH）
* 设计思路：
 */

/*
* 从视频中截取图片
*
 */

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

/*
* 安全的结束ffmpeg进程
 */
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
		//当前小时才刚刚开始，没有文件写入
		log.Println("maxNum:", len(files))
		rec = ImgRootUrl + imgDirNow + "/" + "classify" + strconv.Itoa(len(files)) + ".jpg"
	} else {
		timeLast := timeNow.Add(-time.Hour)
		//查找上一个小时的最新文件返回，这里存在一个BUG,晚上00：00时会触发
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
