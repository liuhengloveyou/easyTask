package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func UpdateTaskAPI(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//r.ParseForm()
	//ttype, name, tid, stat, msg := r.FormValue("type"), r.FormValue("name"), r.FormValue("tid"), r.FormValue("stat"), r.FormValue("msg")
	//if "" == ttype || "" == name || "" == stat || "" == tid {
	//	this.writeErr(w, http.StatusBadRequest, []byte(UPTASKUSAGE))
	//	return
	//}
	//
	//stati, err := strconv.Atoi(stat)
	//if nil != err || (stati != -1 && stati != 1) {
	//	this.writeErr(w, http.StatusBadRequest, []byte("request param stat err"))
	//	glog.Infoln("stat ERR: ", stat)
	//	return
	//}
	//if -1 == stati && "" == msg {
	//	this.writeErr(w, http.StatusBadRequest, []byte("request param msg err"))
	//	glog.Errorln("msg nil: ", stat)
	//	return
	//}
	//if 1 == stati {
	//	stati = 3
	//}
	//
	//taskTypeOne, rapperOne := GetRapper(ttype, name)
	//if taskTypeOne == nil {
	//	this.writeErr(w, http.StatusBadRequest, []byte("no such task type"))
	//	glog.Errorln("putask type err:", ttype)
	//	return
	//} else if rapperOne == nil {
	//	this.writeErr(w, http.StatusBadRequest, []byte("no such rapper"))
	//	glog.Errorln("getask rapper nil:", ttype, name)
	//	return
	//}
	//
	//taskTypeOne.UpTask(rapperOne, int64(stati), tid, msg)
	//glog.Infoln("upTask: ", ttype, name, tid, stat, msg)
	//
	//w.Write([]byte("OK"))

	return
}