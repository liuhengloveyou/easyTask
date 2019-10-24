package api

import (
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/easyTask/services"

	gocommon "github.com/liuhengloveyou/go-common"
)

func UpdateTaskAPI(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("UpdateTaskAPI body ERR: ", err)
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		return
	}

	logger.Debug("UpdateTaskAPI body: ", string(body))

	if string(body) == "" {
		logger.Errorf("UpdateTaskAPI no body")
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		return
	}

	if err = services.UpdateTaskService(body); err != nil {
		logger.Errorf("UpdateTaskAPI service ERR: %v\n", err.Error())
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, "")

	return
}
