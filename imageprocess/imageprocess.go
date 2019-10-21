package imageprocess

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gofaces/dlib_api"
	"gofaces/rtsp"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const dataDir = "./data"
const workerDir = "/home/eagle/gofaces/workers/"

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
		faces, err := faceRec.RecognizeFile(img)
		if err == nil {
			if 1 == len(faces) {
				samples = append(samples, faces[0].Descriptor)
				cats = append(cats, int32(len(samples)))
				labels = append(labels, name)
				ch <- "success"
				faceRec.SetSamples(samples, cats)
				/*将识别到的图片放到对应的文件夹中*/

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

func BuildClassifier() int {
	var err error
	var rec = 0
	if faceRec == nil {
		faceRec, err = dlib_api.NewRecognizer(dataDir)
		if nil == err {
			/*查看workers文件夹下面的所有图片，并对其建模*/
			rd, err := ioutil.ReadDir(workerDir)
			if nil == err {
				for _, fi := range rd {
					image := workerDir + fi.Name()
					face, err := faceRec.RecognizeFile(image)
					if nil == err {
						if 1 == len(face) {
							samples = append(samples, face[0].Descriptor)
							cats = append(cats, int32(len(samples)))
							labels = append(labels, fi.Name()[:len(fi.Name())-4])
							faceRec.SetSamples(samples, cats)
							fmt.Println("get faces:" + fi.Name()[:len(fi.Name())-4])
						}
					} else {
						rec = 1
					}
				}
			} else {
				rec = 2
			}
		} else {
			rec = 3
		}
	}
	return rec
}

func FaceDetections(ch chan<- string) {
	var err error
	var imgOld = ""
	var imgDirOld = ""

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
		if 0 == strings.Compare(img, imgOld) {
			continue
		}
		log.Println("latestImage:" + img)
		imgOld = img
		faces, err := faceRec.RecognizeFile(img)
		if err == nil {
			if 0 != len(faces) {
				results := ""
				/*将识别到的图片放到对应的文件夹中,文件以出现的时间命名（分钟）*/
				imgDir := strconv.Itoa(time.Now().Minute()) + "-" +
					strconv.Itoa(time.Now().Hour()) + "-" +
					strconv.Itoa(time.Now().Day()) + "-" +
					time.Now().Month().String() + "-" +
					strconv.Itoa(time.Now().Year())
				results = results + "{\"time\":\"" + imgDir + "\","
				results = results + "\"num\":" + strconv.Itoa(len(faces)) + ","
				imgDir = rtsp.ModelRoot + imgDir
				if 0 != strings.Compare(imgDir, imgDirOld) {
					/*这里会造成进程灾难*/
					go SaveImages(img, imgDir)
					imgDirOld = imgDir
				}
				log.Println("num of faces" + strconv.Itoa(len(faces)))
				for i := 0; i < len(faces); i++ {
					catsId := faceRec.Classify(faces[i].Descriptor)
					fmt.Println("catsId:" + strconv.Itoa(catsId))
					if catsId <= 0 {
						results = results + "\"who\":\"unknown\","
					} else {
						results = results + "\"who\":\"" + labels[catsId-1] + "\","
					}
				}
				results = results + "\"end\":1}"
				log.Println(results)
				ch <- results
			}
		}
	}
}

func SaveImages(img string, imgDir string) {
	var imgOld = " "
	mkdir := exec.Command("mkdir", "-p", imgDir)
	err := mkdir.Run()
	if err != nil {
		log.Println("mkdir erro:%v", err)
		return
	}
	for i := 0; i < 10; {
		img = rtsp.GetLatestImage()
		if 0 == strings.Compare(img, imgOld) {
			continue
		}
		imgOld = img
		cp := exec.Command("cp", img, imgDir)
		err = cp.Start()
		cp.Wait()
		i++
	}
}
