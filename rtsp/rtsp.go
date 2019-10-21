package rtsp

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//ffmpeg -i "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream" -r 1  -y /home/eagle/rtmp/images/image%d.jpg 每秒截取一张图片
//ffmpeg -y -i "rtmp://58.200.131.2:1935/livetv/hunantv" -ss 00:00:01 -vframes 1 -f image2 /home/eagle/rtmp/images/image1.jpg
//ffmpeg -y -i "rtmp://58.200.131.2:1935/livetv/hunantv" -r 1  -y /home/gofaces/rtmp/images/image%d.jpg
//const RtspUrl = "rtsp://admin:YHDYPD@10.108.219.232:554/h264/ch1/main/av_stream"
const RtspUrl = "rtmp://rtmp01open.ys7.com/openlive/77c99cc69e4443aeaef0c1fd8ac1e5e6.hd"
const ImgRootUrl = "/home/eagle/gofaces/rtmp/"
const ModelRoot = "/home/eagle/gofaces/history/"

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
func VideoCaptureStart(imgName string, imgDir string) (string, string) {
	commandFmt := "ffmpeg  -y -i \"%v\" -r 5 %v"
	img := ImgRootUrl + imgDir
	log.Println("mkdir -p ", img)
	mkdir := exec.Command("mkdir", "-p", img)
	err := mkdir.Run()
	if nil == err {
		img = img + "/" + imgName
		command := fmt.Sprintf(commandFmt, RtspUrl, img)
		log.Println(command)
		cmd := exec.Command("sh", "-c", command)
		err = cmd.Start()
		if nil == err {
			go cmd.Wait()
			return strconv.Itoa(cmd.Process.Pid + 1), "success"
		} else {
			return "0", "error"
		}
		go mkdir.Wait()
	}
	return "0", "error"
}

/*
* 安全的结束ffmpeg进程
 */
func VideoCaptureStop(pid string) string {
	cmd := exec.Command("sh", "-c", "pkill ffmpeg")
	err := cmd.Run()
	if nil == err {
		return "success"
		go cmd.Wait()
	}
	return "error"
}

/*
 用于协调ffmpeg进程
*/
func VideoCaptureHandler1() {
	imgDir := strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Year())
	log.Println(imgDir)
	pid, rec := VideoCaptureStart("classify%d.jpg", imgDir)
	log.Println("restart ffmpeg")
	log.Println("ffmpeg pid:" + pid)
	if 0 == strings.Compare("success", rec) {
		ticker1 := time.NewTicker(time.Minute * time.Duration(60-time.Now().Minute()))
		log.Println("ticker time out", <-ticker1.C)
		rec = VideoCaptureStop(pid)
		ticker1.Stop()
	} else {
		log.Fatalf("ffmpeg 出错，请检查系统")
	}
	ticker := time.NewTicker(time.Hour)
	for true {
		imgDir = strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Year())
		log.Println(imgDir)
		log.Println("restart ffmpeg")
		pid, rec = VideoCaptureStart("classify%d.jpg", imgDir)
		log.Println("ffmpeg pid:" + pid)
		if 0 == strings.Compare("success", rec) {
			log.Println("ffmpeg pid :", pid)
			log.Println("ticker time out", <-ticker.C)
			VideoCaptureStop(pid)
			go CleanOldImages()
		} else {
			log.Fatalf("ffmpeg error please check system!")
		}
	}
}

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
		rec = ImgRootUrl + imgDirNow + "/" + "classify" + strconv.Itoa(len(files)) + ".jpg"
	} else {
		timeLast := timeNow.Add(-time.Hour)
		imgDirLast := strconv.Itoa(timeLast.Hour()) + "-" + strconv.Itoa(timeLast.Day()) + "-" + timeLast.Month().String() + "-" + strconv.Itoa(timeLast.Year())
		files, _ = ioutil.ReadDir(ImgRootUrl + imgDirLast + "/")
		if len(files) > 0 {
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

func VideoCaptureHandler() {
	var pid string
	ch := make(chan string)
	imgDir := strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Year())
	log.Println(imgDir)
	go VideoCaptureStart1("classify%d.jpg", imgDir, ch)
	if 0 == strings.Compare("success", <-ch) {
		ticker1 := time.NewTicker(time.Minute * time.Duration(60-time.Now().Minute()))
		pid = <-ch
		log.Println("ticker time out", <-ticker1.C)
		go VideoCaptureStop1(pid, ch)
		<-ch
		ticker1.Stop()
	} else {
		log.Fatalf("ffmpeg 出错，请检查系统")
	}
	ticker := time.NewTicker(time.Hour)
	for true {
		imgDir = strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Year())
		log.Println(imgDir)
		go VideoCaptureStart1("classify%d.jpg", imgDir, ch)
		if 0 == strings.Compare("success", <-ch) {
			pid = <-ch
			log.Println("ffmpeg pid :", pid)
			log.Println("ticker time out", <-ticker.C)
			go VideoCaptureStop1(pid, ch)
			<-ch
			go CleanOldImages()
		} else {
			log.Fatalf("ffmpeg error please check system!")
		}
	}
}
