package main

import (
	"net/http"

	"github.com/golang/glog"
)

type sayhiHandler struct {}

func (this *sayhiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (this *sayhiHandler) doGet (w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /sayhi?type=typename&name=rappername"
	
	r.ParseForm()
	ttype, name := r.FormValue("type"), r.FormValue("name")
	if ttype == "" || name == "" {
		glog.Errorln("sayhi ERR:", ttype, name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}

	taskTypeOne, rapperOne := GetRapper(ttype, name)
	if taskTypeOne == nil {
		glog.Errorln("sayhi to nil type:", ttype)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such task type."))
		return
	}

	if rapperOne == nil {
		taskTypeOne.rappers[name] = NewRapper()
	} else {
		taskTypeOne.resetRapper(rapperOne)
		rapperOne.Beat(true)
	}

	w.Write([]byte("OK"))
	return
}
