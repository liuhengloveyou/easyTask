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
	if r.Method != "GET" {
		handlePutTask(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}


func handlePutTask(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ttype, rid, info := r.FormValue("type"), r.FormValue("rid"), r.FormValue("info")
	if "" == ttype || "" == rid || "" == info {
		glog.Errorln("putask param nil")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("param err"))
		return
	}

	m := md5.New()
	io.WriteString(m, info)
	taskid := fmt.Sprintf("%x", m.Sum(nil))
	
	err := newTask(ttype, taskid, rid, info)
	if err != nil {
		glog.Errorln("newTask ERR:", err, ttype, rid, info)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("db err"))
		return
	}

	glog.Errorf("DATA putTask: %s %s %s %s", taskid, ttype, rid, info)
	w.Write([]byte(taskid))

	return
}
