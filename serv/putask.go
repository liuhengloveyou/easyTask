package main

import (
	"io"
	"fmt"
	"crypto/md5"
	"net/http"

	"github.com/golang/glog"
)

type putTaskHandler struct {}

func (this *putTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getPutTask(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}


func getPutTask(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /putask?type=typename&rid=recordid&info=taskinfo"

	r.ParseForm()
	ttype, rid, info := r.FormValue("type"), r.FormValue("rid"), r.FormValue("info")
	if "" == ttype || "" == rid || "" == info {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}

	taskTypeOne, ok := TaskTypes[ttype]
	if ok == false {
		glog.Errorln("putask type err:", ttype)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such task type"))
		return
	}
	
	m := md5.New()
	io.WriteString(m, info)
	taskid := fmt.Sprintf("%x", m.Sum(nil))

	taskTypeOne.newTask(&TaskInfo{Tid:taskid, Rid: rid, Info: info})
	
	glog.Errorf("DATA putTask: %s %s %s %s", taskid, ttype, rid, info)
	w.Write([]byte(taskid))

	return
}
