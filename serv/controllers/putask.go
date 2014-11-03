package controllers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	. "easyTask/serv/models"
	. "easyTask/serv/common"
	
	"github.com/golang/glog"
)

type PutTaskHandler struct{}

func (this *PutTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (this *PutTaskHandler) doGet(w http.ResponseWriter, r *http.Request) {
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
	if inSize >= int64(ConfJson["taskBuffSize"].(float64)) {
		this.writeErr(w, http.StatusInternalServerError, []byte("server to busy"))
		glog.Errorln("server to busy err:", inSize, backSize, int(ConfJson["taskBuffSize"].(float64)))
		return
	}

	var stat int64 = 1
	if backSize < int64(ConfJson["taskBuffSize"].(float64)) && taskTypeOne.RapperNum() > 0 {
		stat = 2
	}
	
	m := md5.New()
	io.WriteString(m, info)
	taskid := fmt.Sprintf("%x", m.Sum(nil))
	
	taskTypeOne.NewTask(&TaskInfo{Tid: taskid, Rid: rid, Info: info}, stat)

	glog.Errorf("DATA putTask: %s %s %s %s", taskid, ttype, rid, info)
	w.Write([]byte(taskid))

	return
}

func (this *PutTaskHandler) writeErr(w http.ResponseWriter, statCode int, body []byte) {
	w.WriteHeader(statCode)
	w.Write(body)
}
