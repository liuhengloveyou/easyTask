package api

import (
	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
	"github.com/liuhengloveyou/easyTask/services"
	"net/http"
	"strconv"

	gocommon "github.com/liuhengloveyou/go-common"
)


func QueryTaskAPI(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseForm()

	taskType, name, num := r.FormValue("type"), r.FormValue("name"), r.FormValue("num")
	if "" == taskType || "" == name || "" == num {
		gocommon.HttpErr(w, http.StatusOK, -1, "参数错误")
		return
	}

	inum, err := strconv.Atoi(num)
	if err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		glog.Errorln("QueryTaskAPI num ERR:", num)
		return
	}
	if inum > 10 {
		inum = 10 // 一次最多10个
	}
	glog.Infoln("QueryTaskAPI: ", name, taskType, inum)

	tasks, err := services.QueryTaskService(taskType, inum)
	if err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		logger.Error("QueryTaskAPI ERR: " + err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, tasks)
	logger.Infof("QueryTaskAPI OK: %v %v %v %v", name, taskType, inum, tasks)

	return
}
