package api

import (
	"io/ioutil"
	"net/http"

	"github.com/liuhengloveyou/easyTask/services"
	passportcommon "github.com/liuhengloveyou/passport/dao"
	passport "github.com/liuhengloveyou/passport/face"
	"github.com/liuhengloveyou/passport/sessions"

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

func RedoTaskAPI(w http.ResponseWriter, r *http.Request) {
	var UserID int64
	if r.Context().Value("session") != nil {
		UserID = r.Context().Value("session").(*sessions.Session).Values[passport.SessUserInfoKey].(passportcommon.User).UID
	}
	if UserID <= 0 {
		gocommon.HttpErr(w, http.StatusForbidden, -1, "")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("RedoTaskAPI body ERR: ", err)
		gocommon.HttpErr(w, http.StatusOK, -1, "请求错误")
		return
	}

	logger.Debug("RedoTaskAPI body: ", string(body))

	if len(body) < 5 {
		logger.Errorf("RedoTaskAPI no body")
		gocommon.HttpErr(w, http.StatusOK, -1, "请求错误")
		return
	}

	if err = services.RedoTaskService(UserID, body); err != nil {
		logger.Errorf("RedoTaskAPI service ERR: %v\n", err.Error())
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, "")

	return

}
