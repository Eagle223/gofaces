package rtsp

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"time"
)

//ffmpeg -i "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream" -r 1  -y /home/eagle/rtmp/images/image%d.jpg 每秒截取一张图片
//ffmpeg -y -i "rtmp://58.200.131.2:1935/livetv/hunantv" -ss 00:00:01 -vframes 1 -f image2 /home/eagle/rtmp/images/image1.jpg

const RtspUrl = "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream"
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
	imgDir := strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day()) + "-" + time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Year())

	files, err := ioutil.ReadDir(ImgRootUrl + imgDir + "/")
	sort.Slice(files, func(i, j int) bool {

		return false
	})
}
