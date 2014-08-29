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
		this.writeErr(w, http.StatusBadRequest, []byte(USAGE))
		return
	}

	stati, err := strconv.Atoi(stat)
	if nil != err || (stati != -1 && stati != 1) {
		this.writeErr(w, http.StatusBadRequest, []byte("request param stat err"))
		glog.Infoln("stat ERR: ", stat)
		return
	}
	if -1 == stati && "" == msg {
		this.writeErr(w, http.StatusBadRequest, []byte("request param msg err"))
		glog.Errorln("msg nil: ", stat)
		return
	}
	if 1 == stati {
		stati = 3
	}

	taskTypeOne, rapperOne := GetRapper(ttype, name)
	if taskTypeOne == nil {
		this.writeErr(w, http.StatusBadRequest, []byte("no such task type"))
		glog.Errorln("putask type err:", ttype)
		return
	} else if rapperOne == nil {
		this.writeErr(w, http.StatusBadRequest, []byte("no such rapper"))
		glog.Errorln("getask rapper nil:", ttype, name)
		return
	}
	
	taskTypeOne.upTask(rapperOne, int64(stati), tid, msg)
	glog.Infoln("upTask: ", ttype, name, tid, stat, msg)
		
	w.Write([]byte("OK"))

	return
}

func (this *upTaskHandler) writeErr(w http.ResponseWriter, statCode int, body []byte) {
	w.WriteHeader(statCode)
	w.Write(body)
}
