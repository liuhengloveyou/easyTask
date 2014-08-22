package main

import (
	"net/http"

	"github.com/golang/glog"
)

type handleSayhi struct {}

func (this *handleSayhi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getSayhi(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func getSayhi(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /sayhi?type=typename&name=rappername"
	
	r.ParseForm()
	ttype, name := r.FormValue("type"), r.FormValue("name")
	if ttype == "" || name == "" {
		glog.Errorln("sayhi ERR:", ttype, name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}
	
	taskTypeOne, ok := TaskTypes[ttype]
	if ok != true {
		glog.Errorln("sayhi to nil type:", ttype)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such task type."))
		return
	}

	rapperOne, ok := taskTypeOne.rappers[name]
	if ok == false {
		taskTypeOne.rappers[name] = NewRapper()
	} else {
		taskTypeOne.resetRapper(rapperOne)
	}

	w.Write([]byte("OK"))
	return
}
