package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/glog"
)

type putTaskHandler struct{}

func (this *putTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *putTaskHandler) doGet(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /putask?type=typename&rid=recordid&info=taskinfo"

	r.ParseForm()
	ttype, rid, info := r.FormValue("type"), r.FormValue("rid"), r.FormValue("info")
	if "" == ttype || "" == rid || "" == info {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}

	taskTypeOne, _ := GetRapper(ttype, "")
	if taskTypeOne == nil {
		this.writeErr(w, http.StatusBadRequest, []byte("no such task type"))
		glog.Errorln("putask type err:", ttype)
		return
	}

	inSize, backSize := taskTypeOne.BuffSize()
	if inSize >= int64(confJson["taskBuffSize"].(float64)) {
		this.writeErr(w, http.StatusInternalServerError, []byte("server to busy"))
		glog.Errorln("server to busy err:", inSize, backSize, int(confJson["taskBuffSize"].(float64)))
		return
	}

	var stat int64 = 1
	if backSize < int64(confJson["taskBuffSize"].(float64)) && len(taskTypeOne.rappers) > 0 {
		stat = 2
	}
	
	m := md5.New()
	io.WriteString(m, info)
	taskid := fmt.Sprintf("%x", m.Sum(nil))
	
	taskTypeOne.newTask(&TaskInfo{Tid: taskid, Rid: rid, Info: info}, stat)

	glog.Errorf("DATA putTask: %s %s %s %s", taskid, ttype, rid, info)
	w.Write([]byte(taskid))

	return
}

func (this *putTaskHandler) writeErr(w http.ResponseWriter, statCode int, body []byte) {
	w.WriteHeader(statCode)
	w.Write(body)
}
