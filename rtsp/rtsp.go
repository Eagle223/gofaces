package rtsp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os/exec"
	"strings"
)

//ffmpeg -i "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream" -r 1 /home/eagle/rtmp/images/image%d.jpg 每秒截取一张图片
//ffmpeg -y -i "rtmp://58.200.131.2:1935/livetv/hunantv" -ss 00:00:01 -vframes 1 -f image2 /home/eagle/rtmp/images/image1.jpg

const RtspUrl = "rtsp://admin:YHDYPD@192.168.184.180:554/h264/ch1/main/av_stream"
const ImgRootUrl = "/home/eagle/rtmp/"

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

func VideoCaptureStart1(imgNames string, ch chan<- string) {
	commandFmt := "ffmpeg -i \"%v\" -r 1 %v"
	img := ImgRootUrl + imgNames
	mkdir := exec.Command("mkdir", "-p", img)
	err := mkdir.Run()
	if err != nil {
		ch <- fmt.Sprintf("mkdir erro:%v", err)
		return
	}
	img = img + "/image%d.jpg"
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
	cmd.Wait()
	mkdir.Wait()
}

/*
* 安全的结束ffmpeg进程
 */
func VideoCaptureStop1(ch chan<- string) {
	cmd := exec.Command("sh", "-c", "pkill ffmpeg")
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

func VideoCaptureStart(c *gin.Context) {
	imgName := "/images/"
	ch := make(chan string)
	go VideoCaptureStart1(imgName, ch)
	rec := <-ch
	log.Println(rec)
	if 0 == strings.Compare("success", rec) {
		c.JSON(200, gin.H{
			"message": "操作完成",
			"rtspUrl": RtspUrl,
		})
	} else {
		c.JSON(500, gin.H{
			"message": rec,
		})
	}
}

func VideoCaptureStop(c *gin.Context) {
	ch := make(chan string)
	go VideoCaptureStop1(ch)
	rec := <-ch
	if 0 == strings.Compare("success", rec) {
		c.JSON(200, gin.H{
			"message": "关闭成功",
		})
	} else {
		c.JSON(500, gin.H{
			"message": rec,
		})
	}
}
