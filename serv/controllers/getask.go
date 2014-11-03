package controllers

import (
	"strconv"
	"net/http"
	"encoding/json"

	. "easyTask/serv/models"
	. "easyTask/serv/common"
	
	"github.com/golang/glog"
)

type GetTaskHandler struct {}

func (this *GetTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if Sig != "" {
		this.writeErr(w, http.StatusServiceUnavailable, []byte(Sig))
		return
	}
	
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *GetTaskHandler) doGet(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /getask?type=typename&name=rappername&num=123"
	
	r.ParseForm()
	ttype, name, num := r.FormValue("type"), r.FormValue("name"), r.FormValue("num")
	if "" == ttype || "" == name || "" == num {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}
	inum, err := strconv.Atoi(num)
	if err != nil {
		this.writeErr(w, http.StatusBadRequest, []byte("num err"))
		glog.Errorln("getask num ERR:", num)
		return
	}
	if inum >= int(ConfJson["MaxTaskPerRapper"].(float64)) {
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

	if rapperOne.TaskSize() > int(ConfJson["MaxTaskPerRapper"].(float64)) {
		this.writeErr(w, http.StatusBadRequest, []byte("to much tasks"))
		glog.Errorln("getask to mach ERR:", rapperOne.TaskSize, int(ConfJson["MaxTaskPerRapper"].(float64)))
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
	w.Write(jsonByte)
	glog.Infoln("getask OK: ", string(jsonByte))
	
	return
}

func (this *GetTaskHandler) writeErr(w http.ResponseWriter, statCode int, body []byte) {
	w.WriteHeader(statCode)
	w.Write(body)
}
