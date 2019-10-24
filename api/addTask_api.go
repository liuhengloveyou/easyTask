package api

import (
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/easyTask/services"

	gocommon "github.com/liuhengloveyou/go-common"
	passportcommon "github.com/liuhengloveyou/passport/dao"
	passport "github.com/liuhengloveyou/passport/face"
	"github.com/liuhengloveyou/passport/sessions"
)

func AddTaskAPI(w http.ResponseWriter, r *http.Request) {
	var UserID int64
	if r.Context().Value("session") != nil {
		UserID = r.Context().Value("session").(*sessions.Session).Values[passport.SessUserInfoKey].(passportcommon.User).UID
	}

	taskType := r.FormValue("type")
	if taskType == "" {
		logger.Errorf("AddTask type para ERR")
		gocommon.HttpErr(w, http.StatusOK, -1, "参数错误")
		return
	}

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

	id, err := services.AddTaskService(UserID, taskType, string(body), false)
	if err != nil {
		logger.Errorf("AddTask service ERR: %v\n", err.Error())
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, id)
	logger.Infof("AddTask OK: %#v.\n", id, body)

	return
}

func AddTaskBatchAPI(w http.ResponseWriter, r *http.Request) {
	var UserID int64
	if r.Context().Value("session") != nil {
		UserID = r.Context().Value("session").(*sessions.Session).Values[passport.SessUserInfoKey].(passportcommon.User).UID
	}

	taskType := r.FormValue("type")
	if taskType == "" {
		logger.Errorf("AddTask type para ERR")
		gocommon.HttpErr(w, http.StatusOK, -1, "参数错误")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("AddTask body ERR: ", err)
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		return
	}

	logger.Debug("AddTask: ", UserID, taskType, string(body))

	if string(body) == "" {
		logger.Errorf("AddTask requst ERR")
		gocommon.HttpErr(w, http.StatusOK, -1, "读取数据出错")
		return
	}

	ids, err := services.AddTaskService(UserID, taskType, string(body), true)
	if err != nil {
		logger.Errorf("AddTask service ERR: %v\n", err.Error())
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, ids)
	logger.Infof("AddTask OK: %#v\n", ids)

	return
}
