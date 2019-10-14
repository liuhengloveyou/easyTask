package api

import (
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/easyTask/services"

	"github.com/julienschmidt/httprouter"
	gocommon "github.com/liuhengloveyou/go-common"
)

func AddTaskAPI(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("AddTask body ERR: ", err)
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		return
	}

	logger.Debug("AddTask body: ", string(body))

	if string(body) == "" {
		logger.Errorf("AddTask requst ERR")
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		return
	}

	id, err := services.AddTaskService(body)
	if err != nil {
		logger.Errorf("AddTask service ERR: %v\n", err.Error())
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, id)
	logger.Infof("AddTask OK: %#v.\n", id, body)

	return
}