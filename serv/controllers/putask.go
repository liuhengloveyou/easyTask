package controllers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	. "github.com/liuhengloveyou/easyTask/common"
	. "github.com/liuhengloveyou/easyTask/serv/models"

	"github.com/golang/glog"
)

const PUTTASKUSAGE = "GET /putask?type=typename&rid=recordid&info=taskinfo"

type PutTaskHandler struct{}

func (this *PutTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		this.doGet(w, r)
	} else {
		this.writeErr(w, http.StatusMethodNotAllowed, []byte(PUTTASKUSAGE))
	}

	glog.Flush()
	return
}

func (this *PutTaskHandler) doGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ttype, rid, info := r.FormValue("type"), r.FormValue("rid"), r.FormValue("info")
	if "" == ttype || "" == rid || "" == info {
		this.writeErr(w, http.StatusBadRequest, []byte(PUTTASKUSAGE))
		return
	}

	taskTypeOne, _ := GetRapper(ttype, "")
	if taskTypeOne == nil {
		this.writeErr(w, http.StatusBadRequest, []byte("no such task type"))
		glog.Errorln("putask type err:", ttype)
		return
	}

	inSize, outSize := taskTypeOne.BuffSize()
	if inSize >= int64(ConfJson["inBuffSize"].(float64)) {
		this.writeErr(w, http.StatusInternalServerError, []byte("service to busy"))
		glog.Errorln("server to busy err:", inSize, outSize, int(ConfJson["inBuffSize"].(float64)))
		return
	}

	var stat int64 = 1
	if outSize < int64(ConfJson["outBuffSize"].(float64)) && taskTypeOne.RapperNum() > 0 {
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
