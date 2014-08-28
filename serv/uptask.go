package main

import (
	"strconv"
	"net/http"

	"github.com/golang/glog"
)

type upTaskHandler struct {}

func (this *upTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *upTaskHandler) doGet(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /uptask?type=typename&name=rappername&tid=taskid&stat=1|-1&msg=errormsg"
	
	r.ParseForm()
	ttype, name, tid, stat, msg := r.FormValue("type"), r.FormValue("name"), r.FormValue("tid"), r.FormValue("stat"), r.FormValue("msg")
	if "" == ttype || "" == name || "" == stat || "" == tid {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}

	stati, err := strconv.Atoi(stat)
	if nil != err || (stati != -1 && stati != 1) {
		glog.Infoln("stat ERR: ", stat)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request param stat err"))
		return
	}
	if -1 == stati && "" == msg {
		glog.Infoln("msg nil: ", stat)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request param msg err"))
		return
	}
	if 1 == stati {
		stati = 3
	}

	taskTypeOne, rapperOne := GetRapper(ttype, name)
	if taskTypeOne == nil {
		glog.Errorln("putask type err:", ttype)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such task type"))
		return
	} else if rapperOne == nil {
		glog.Errorln("getask rapper nil:", ttype, name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such rapper"))
		return
	}
	
	taskTypeOne.upTask(rapperOne, int64(stati), tid, msg)
	glog.Infoln("upTask: ", ttype, name, tid, stat, msg)
		
	w.Write([]byte("OK"))

	return
}

