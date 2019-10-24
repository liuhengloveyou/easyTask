package api

import (
	"net/http"
	"strconv"

	"github.com/liuhengloveyou/easyTask/services"

	gocommon "github.com/liuhengloveyou/go-common"
)

func GetTaskAPI(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	taskType, name, num := r.FormValue("type"), r.FormValue("name"), r.FormValue("num")
	if "" == taskType || "" == num {
		gocommon.HttpErr(w, http.StatusOK, -1, "参数错误")
		return
	}

	if "" == name {
		gocommon.HttpErr(w, http.StatusOK, -1, "参数错误")
		return
	}

	inum, err := strconv.Atoi(num)
	if err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		logger.Errorf("QueryTaskAPI num ERR: %v\n", num)
		return
	}
	if inum > 10 {
		inum = 10 // 一次最多10个
	}
	logger.Infof("QueryTaskAPI: %v %v %v\n", name, taskType, inum)

	tasks, err := services.GetTaskService(taskType, inum)
	if err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		logger.Errorf("QueryTaskAPI ERR: %v %v %v\n", name, taskType, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, tasks)
	logger.Infof("QueryTaskAPI OK: %v %v %v\n", name, taskType, inum)

	return
}
