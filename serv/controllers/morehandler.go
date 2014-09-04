package controllers

import (
	"net/http"

	. "easyTask/serv/models"
	
	"github.com/golang/glog"
)

func HandleNewTaskType(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /newtype?name=typename"

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return

	} else if len(name) > 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name's len < 8"))
		return
	}

	err := CreateDB(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	TaskTypes[name] = NewTaskType()
	TaskTypes[name].Name = name

	w.Write([]byte("OK"))

	return
}

func HandleBeat(w http.ResponseWriter, r *http.Request) {
	const USAGE = "GET /beat?type=tasktype&name=rappername"

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	ttype, name := r.FormValue("type"), r.FormValue("name")
	if ttype == "" || name == "" {
		glog.Errorln("beat ERR:", ttype, name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(USAGE))
		return
	}

	taskTypeOne, rapperOne := GetRapper(ttype, name)
	if taskTypeOne == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such task type."))
		return
	} else if rapperOne == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no such rapper"))
		return
	}

	rapperOne.Beat(true)

	w.Write([]byte("OK"))
	return
}
