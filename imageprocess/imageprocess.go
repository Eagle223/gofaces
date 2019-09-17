package imageprocess

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gofaces/dlib_api"
	"gofaces/rtsp"
	"log"
	"path/filepath"
	"strings"
	"sync/atomic"
)

const dataDir = "./data"

var samples []dlib_api.Descriptor
var cats []int32
var labels []string
var facecats int32

func init() {
	facecats = 0
}

/*
 * 换一个思路
 */

func buildFaceModle(imgDir string, ch chan<- string) {
	rec, err := dlib_api.NewRecognizer(dataDir)
	if err != nil {
		log.Println("Can't init face recognizer: %v", err)
	}
	defer rec.Close()
	i := 1
	for true {
		img := fmt.Sprintf(rtsp.ImgRootUrl+imgDir+"/image%d.jpg", i)
		fmt.Println(img)
		faces, err := rec.RecognizeFile(img)
		if err == nil {
			if 1 == len(faces) {
				ch <- "success"
				//这里需要加锁
				atomic.AddInt32(&facecats, 1)
				samples = append(samples, faces[0].Descriptor)
				labels = append(labels, imgDir)
				cats = append(cats, int32(facecats))
				break
			}
			i++
		}
	}
}
func classifyFace(imgDir string, ch chan<- string) {
	// Init the recognizer.
	rec, err := dlib_api.NewRecognizer(dataDir)
	if err != nil {
		log.Println("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()
	rec.SetSamples(samples, cats)
	testImage := filepath.Join(rtsp.ImgRootUrl, imgDir)

	for i := 1; i < 50; {
		img := fmt.Sprintf(testImage+"/image%d.jpg", i)
		face, err := rec.RecognizeFile(img)
		if err == nil {
			if 1 == len(face) {
				catID := rec.Classify(face[0].Descriptor)
				if catID >= 0 {
					fmt.Println("catID", catID)
					ch <- labels[catID-1]
					return
				}
			} else {
				log.Println("face != 1")
			}
			i++
		} else {
			log.Println("RecognizeFile error:%v", err)
		}
	}
	ch <- "50"
}

/*
* 为摄像机前的人脸建立识别模型 传入人脸名字
*
 */
func BuildFaceModle(c *gin.Context) {
	faceName := c.Query("facename")
	ch1 := make(chan string)
	ch2 := make(chan string)
	go rtsp.VideoCaptureStart1(faceName, ch1)
	rec := <-ch1
	if 0 == strings.Compare("success", rec) {
		go buildFaceModle(faceName, ch2)
		rec = <-ch2
		if 0 == strings.Compare("success", rec) {
			c.JSON(200, gin.H{
				"message": "构建人脸模型成功",
			})
		} else {
			c.JSON(500, gin.H{
				"message": rec,
			})
		}
		go rtsp.VideoCaptureStop1(ch1)
		fmt.Println("stop capture", <-ch1)
	} else {
		c.JSON(500, gin.H{
			"message": rec,
		})
	}
}

/*
 * 识别摄像机前的人脸
 */
func ClassifyFace(c *gin.Context) {
	ch1 := make(chan string)
	ch2 := make(chan string)
	go rtsp.VideoCaptureStart1("classify", ch1)
	rec := <-ch1
	if 0 == strings.Compare("success", rec) {
		go classifyFace("classify", ch2)
		rec = <-ch2
		go rtsp.VideoCaptureStop1(ch1)
		fmt.Println("stop capture", <-ch1)
		if 0 == strings.Compare("50", rec) {
			c.JSON(404, gin.H{
				"message": "未识别到",
			})
		} else {
			c.JSON(200, gin.H{
				"message": rec,
			})
		}
	} else {
		c.JSON(500, gin.H{
			"message": rec,
		})
	}
	fmt.Println("befor ClassifyFace return")
}
