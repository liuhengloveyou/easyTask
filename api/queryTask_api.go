package api

import (
	"net"
	"net/http"
	"strconv"

	"github.com/liuhengloveyou/easyTask/services"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
	gocommon "github.com/liuhengloveyou/go-common"
)

func QueryTaskAPI(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseForm()

	taskType, name, num := r.FormValue("type"), r.FormValue("name"), r.FormValue("num")
	if "" == taskType || "" == num {
		gocommon.HttpErr(w, http.StatusOK, -1, "参数错误")
		return
	}

	if "" == name {
		name, _, _ = net.SplitHostPort(r.RemoteAddr)
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
	logger.Infof("QueryTaskAPI OK: %v %v %v\n", name, taskType, inum)

	return
}
