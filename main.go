package main

import (
	"encoding/json"
	"gofaces/communication"
	"gofaces/imageprocess"
	"gofaces/rtsp"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	var timeOld = ""
	ch := make(chan map[string]string)
	aliveList := communication.NewAliveList()
	go communication.ServerStart(aliveList)
	go rtsp.VideoCaptureHandler1()
	if 3 == imageprocess.BuildClassifier() {
		log.Println("加载Data 错误！")
		return
	}
	go imageprocess.FaceDetections(ch)
	for true {
		context := <-ch
		if 0 != strings.Compare(context["time"], timeOld) {
			timeOld = context["time"]
			id := strconv.Itoa(rand.Intn(65536))
			contextJson, _ := json.Marshal(context)
			contextString := string(contextJson)
			log.Println("main:", contextString)
			message := communication.Message{
				ID:      id,
				Content: contextString,
				SentAt:  time.Now().Unix(),
				Type:    0,
			}
			aliveList.Broadcast(message)
		} else {
			log.Println("已经发送了信息！")
		}
	}

}
