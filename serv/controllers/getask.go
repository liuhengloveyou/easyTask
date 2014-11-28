package controllers

import (
	"strconv"
	"net/http"
	"encoding/json"

	. "easyTask/serv/models"
	. "easyTask/serv/common"
	
	"github.com/golang/glog"
)

const GETTASKUSAGE = "GET /getask?type=typename&name=rappername&num=123"

type GetTaskHandler struct {}

func (this *GetTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if Sig != "" {
		this.writeErr(w, http.StatusServiceUnavailable, []byte(Sig))
		return
	}
	
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		this.writeErr(w, http.StatusMethodNotAllowed, []byte(GETTASKUSAGE))
	}

	glog.Flush()
	return
}

func (this *GetTaskHandler) doGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ttype, name, num := r.FormValue("type"), r.FormValue("name"), r.FormValue("num")
	if "" == ttype || "" == name || "" == num {
		this.writeErr(w, http.StatusBadRequest, []byte(GETTASKUSAGE))
		return
	}
	
	inum, err := strconv.Atoi(num)
	if err != nil {
		this.writeErr(w, http.StatusBadRequest, []byte("num err"))
		glog.Errorln("getask num ERR:", num)
		return
	}
	if inum >= int(ConfJson["maxTaskPerRapper"].(float64)) {
		this.writeErr(w, http.StatusBadRequest, []byte("num to bag"))
		glog.Errorln("getask num ERR:", num)
		return
	}
	glog.Infoln("getask: ", name, ttype, num)

	taskTypeOne, rapperOne := GetRapper(ttype, name)
	if taskTypeOne == nil {
		this.writeErr(w, http.StatusBadRequest, []byte("no such task type"))
		glog.Errorln("getask type nil:", ttype)
		return
	} else if rapperOne == nil {
		this.writeErr(w, http.StatusBadRequest, []byte("no such rapper"))
		glog.Errorln("getask rapper nil:", ttype, name)
		return
	}

	if rapperOne.TaskSize() > int(ConfJson["maxTaskPerRapper"].(float64)) {
		this.writeErr(w, http.StatusBadRequest, []byte("to much tasks"))
		glog.Errorln("getask to mach ERR:", rapperOne.TaskSize, int(ConfJson["maxTaskPerRapper"].(float64)))
		return
	}
	
	rst := make([]TaskInfo, 0)
	tasks := taskTypeOne.DistTask(rapperOne, inum)
	for _, tn := range tasks {
		if tn != nil {
			rst = append(rst, *tn)
		}
	}
	
	jsonByte, _ := json.Marshal(rst)
	glog.Infoln("getask OK: ", string(jsonByte))
	w.Write(jsonByte)
	
	return
}

func (this *GetTaskHandler) writeErr(w http.ResponseWriter, statCode int, body []byte) {
	w.WriteHeader(statCode)
	w.Write(body)
}
