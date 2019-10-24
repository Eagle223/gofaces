package main

import (
	"encoding/json"
	"gofaces/communication"
	"gofaces/imageprocess"
	"gofaces/rtsp"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
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
	}

}
