package imageprocess

import (
	"fmt"
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

func FaceDetections(ch chan<- map[string]string) {
	var err error
	var imgOld = ""
	var imgDirOld = ""
	result := make(map[string]string)
	if faceRec == nil {
		faceRec, err = dlib_api.NewRecognizer(dataDir)
		if err != nil {
			log.Println("Can't init face recognizer: %v", err)
			result["error"] = "error"
			ch <- result
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
				/*将识别到的图片放到对应的文件夹中,文件以出现的时间命名（分钟）*/
				imgDir := strconv.Itoa(time.Now().Minute()) + "-" +
					strconv.Itoa(time.Now().Hour()) + "-" +
					strconv.Itoa(time.Now().Day()) + "-" +
					time.Now().Month().String() + "-" +
					strconv.Itoa(time.Now().Year())

				result["time"] = imgDir
				result["num"] = strconv.Itoa(len(faces))
				log.Println("num of faces" + strconv.Itoa(len(faces)))
				for i := 0; i < len(faces); i++ {
					catsId := faceRec.Classify(faces[i].Descriptor)
					fmt.Println("catsId:" + strconv.Itoa(catsId))
					if catsId <= 0 {
						result["who"+strconv.Itoa(i)] = "unknown"
					} else {
						result["who"+strconv.Itoa(i)] = labels[catsId-1]
					}
					imgDir = imgDir + "_" + result["who"+strconv.Itoa(i)]

					if 0 != strings.Compare(imgDir, imgDirOld) {
						imgDirOld = imgDir
						go SaveImages(img, imgDir)
					}
				}
				ch <- result
			}
		}
	}
}

func SaveImages(img string, imgDir string) {
	mkdir := exec.Command("mkdir", "-p", imgDir)
	err := mkdir.Run()
	if err != nil {
		log.Println("mkdir erro:%v", err)
		return
	}
	// /home/eagle/gofaces/rtmp/21-22-October-2019/classify137.jpg
	imgNum, _ := strconv.Atoi(img[52 : len(img)-4])
	for i := 0; i < 600; {
		imgToCopy := img[:52] + strconv.Itoa(imgNum+i) + ".jpg"
		cp := exec.Command("cp", imgToCopy, imgDir+"/image"+strconv.Itoa(i)+".jpg")
		err = cp.Run()
		if err != nil {
			continue
		}
		cp.Wait()
		log.Printf("copy %v to %v", imgToCopy, imgDir)
		i++
	}
	/*将保存的文件转换为视频*/
	log.Println("开始转换")
	rtsp.BuildMp4FromImage(imgDir)

	/*这里还有一个BUG,会删除失败*/
	//cmd := exec.Command("rm", "-rf", imgDir)
	//err = cmd.Run()
	//if err != nil{
	//	log.Println("删除历史失败 imgDir:",imgDir)
	//	return
	//}
	//log.Println("删除历史完成 imgDir:", imgDir)
	//cmd.Wait()
}
