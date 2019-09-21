package imageprocess

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gofaces/dlib_api"
	"gofaces/rtsp"
	"log"
	"strings"
)

const dataDir = "./data"

var cats []int32
var labels []string
var samples []dlib_api.Descriptor
var faceRec *dlib_api.Recognizer

func buildFaceModle(name string, ch chan<- string) {
	var err error
	if faceRec == nil {
		faceRec, err = dlib_api.NewRecognizer(dataDir)
		if err != nil {
			log.Println("Can't init face recognizer: %v", err)
			ch <- "error"
			return
		}
	}
	for true {
		img := rtsp.GetLatestImage()
		fmt.Println(img)
		faces, err := faceRec.RecognizeFile(img)
		if err == nil {
			if 1 == len(faces) {
				samples = append(samples, faces[0].Descriptor)
				cats = append(cats, int32(len(samples)))
				labels = append(labels, name)
				ch <- "success"
				faceRec.SetSamples(samples, cats)
			}
		}
	}
}

func classifyFace(ch chan<- string) {
	var catsId = -1
	for i := 1; i < 50; {
		img := rtsp.GetLatestImage()
		face, err := faceRec.RecognizeFile(img)
		if err == nil {
			if 1 == len(face) {
				catsId = faceRec.Classify(face[0].Descriptor)
				if catsId > 0 {
					ch <- labels[catsId]
					break
				}
			}
			i++
		}
	}
	if catsId < 0 {
		ch <- "error"
	}
}

func BuildFaceModle(c *gin.Context) {
	faceName := c.Query("facename")
	ch := make(chan string)
	go buildFaceModle(faceName, ch)
	rec := <-ch
	if 0 == strings.Compare("success", rec) {
		c.JSON(200, gin.H{
			"message": "构建人脸模型成功",
		})
	} else {
		c.JSON(500, gin.H{
			"message": rec,
		})
	}
}

func ClassifyFace(c *gin.Context) {
	ch := make(chan string)
	go classifyFace(ch)
	rec := <-ch
	if 0 == strings.Compare("50", rec) {
		c.JSON(404, gin.H{
			"message": "未识别到",
		})
	} else {
		c.JSON(200, gin.H{
			"message": rec,
		})
	}
}
