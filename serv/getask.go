package main

import (
	"strconv"
	"net/http"
	"encoding/json"
	
	"github.com/golang/glog"
)

type getTaskHandler struct {}

func (this *getTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *getTaskHandler) doGet(w http.ResponseWriter, r *http.Request) {
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
		glog.Errorln("getask num ERR:", num)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("num err"))
		return
	}
	if inum >= int(confJson["MaxTaskPerRapper"].(float64)) {
		glog.Errorln("getask num ERR:", num)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("num to bag"))
		return
	}
	glog.Infoln("getask: ", name, ttype, num)

	taskTypeOne, ok := TaskTypes[ttype]
	if ok == false {
		glog.Errorln("getask type nil:", ttype)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such task type"))
		return
	}

	rapperOne, ok := taskTypeOne.rappers[name]
	if ok == false {
		glog.Errorln("getask rapper nil:", ttype, name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such rapper"))
		return
	}

	if rapperOne.TaskSize() > int(confJson["MaxTaskPerRapper"].(float64)) {
		glog.Errorln("getask to mach ERR:", rapperOne.TaskSize, int(confJson["MaxTaskPerRapper"].(float64)))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("to much tasks"))
		return
	}
	
	rst := make([]TaskInfo, 0)
	tasks := taskTypeOne.distTask(rapperOne, inum)
	for _, tn := range tasks {
		if tn != nil {
			rst = append(rst, *tn)
		}
	}

	jsonByte, _ := json.Marshal(rst)
	w.Write(jsonByte)
	
	return
}
