package controller

import (
	"encoding/json"
	"gofaces/rtsp"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func GetHistory(writer http.ResponseWriter, request *http.Request) {
	/*查看history文件夹下面的所有.mp4，并对其建模*/
	rec := make(map[string]string)
	rd, err := ioutil.ReadDir(rtsp.ModelRoot)
	if nil == err {
		msg := make(map[string]string)
		i := 0
		for _, fi := range rd {
			vedio := fi.Name()
			if 0 == strings.Compare(vedio[len(vedio)-4:], ".mp4") {
				vedioPath := fi.Name()
				msg["vedio"+strconv.Itoa(i)] = "/history/" + vedioPath
				i++
			}
		}
		msg["size"] = strconv.Itoa(i)
		contextJson, _ := json.Marshal(msg)
		contextString := string(contextJson)
		rec["status"] = "200"
		rec["msg"] = contextString
		recJson, _ := json.Marshal(rec)
		writer.Write([]byte(string(recJson)))
	} else {
		rec["status"] = "500"
		rec["msg"] = "error"
		recJson, _ := json.Marshal(rec)
		writer.Write([]byte(string(recJson)))
	}
}
