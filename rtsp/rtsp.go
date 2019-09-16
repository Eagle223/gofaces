package rtsp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

//ffmpeg -i "rtmp://58.200.131.2:1935/livetv/hunantv" -r 1 /home/eagle/rtmp/images/image%d.jpg 每秒截取一张图片
//ffmpeg -y -i "rtmp://58.200.131.2:1935/livetv/hunantv" -ss 00:00:01 -vframes 1 -f image2 /home/eagle/rtmp/images/image1.jpg

const rtspUrl = "rtmp://58.200.131.2:1935/livetv/hunantv"
const imgUrl = "/home/eagle/rtmp/images/"

var cmd *exec.Cmd

/*
* ffmpeg 参数设置
* -i 输入流或者文件
* -r 设置帧数Hz为单位
* -s 设置帧大小（WxH）
* 设计思路：
 */

/*
* 从视频中截取一帧图片
 */
func VideoCapture(vedioUrl string) string {

	return vedioUrl
}

func VideoCaptureStart(c *gin.Context) {
	commandFmt := "ffmpeg -i \"%v\" -r 1 %v"
	img := imgUrl + "image%d.jpg"
	command := fmt.Sprintf(commandFmt, rtspUrl, img)
	fmt.Println("command :" + command)
	cmd = exec.Command("sh", "-c", command)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("err : %v", err)
	}
	testImagePristin := filepath.Join(imgUrl + "image1.jpg")
	fd, err := os.Open(testImagePristin)
	for err != nil {
		fd, err = os.Open(testImagePristin)
	}
	defer fd.Close()

	c.JSON(200, gin.H{
		"message": "启动视频截图成功！",
		"rtspUrl": rtspUrl,
	})
}

func VideoCaptureStop(c *gin.Context) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	c.JSON(200, gin.H{
		"message": "关闭视频截图成功！",
		"pid":     cmd.Process.Pid,
	})
}
