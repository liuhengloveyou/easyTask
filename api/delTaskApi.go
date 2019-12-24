package api

import (
	"net/http"
	"strconv"

	"github.com/liuhengloveyou/easyTask/services"

	gocommon "github.com/liuhengloveyou/go-common"
	passportcommon "github.com/liuhengloveyou/passport/dao"
	passport "github.com/liuhengloveyou/passport/face"
	"github.com/liuhengloveyou/passport/sessions"
)

func DeleteTaskAPI(w http.ResponseWriter, r *http.Request) {
	var UserID int64
	if r.Context().Value("session") != nil {
		UserID = r.Context().Value("session").(*sessions.Session).Values[passport.SessUserInfoKey].(passportcommon.User).UID
	}
	if UserID <= 0 {
		gocommon.HttpErr(w, http.StatusForbidden, -1, "")
		return
	}

	taskID, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if taskID <= 0 {
		logger.Errorf("DeleteTaskAPI id nil")
		gocommon.HttpErr(w, http.StatusOK, -1, "请求错误")
		return
	}

	logger.Debug("DeleteTaskAPI: ", taskID, UserID)

	err := services.DeleteTaskService(taskID, UserID)
	if err != nil {
		logger.Errorf("DeleteTaskAPI service ERR: %v\n", err.Error())
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, "")
	logger.Infof("DeleteTaskAPI OK: %#v.\n", taskID, UserID)

	return
}
