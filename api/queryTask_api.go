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

func QueryTaskAPI(w http.ResponseWriter, r *http.Request) {
	var UserID int64
	if r.Context().Value("session") != nil {
		UserID = r.Context().Value("session").(*sessions.Session).Values[passport.SessUserInfoKey].(passportcommon.User).UID
	}

	r.ParseForm()

	taskType := r.FormValue("type")
	pageNO, _ := strconv.ParseUint(r.FormValue("pageno"), 10, 64)
	pageSize, _ := strconv.ParseUint(r.FormValue("pagesize"), 10, 64)

	if pageNO < 1 {
		pageNO = 1
	}
	if pageSize > 100 {
		pageSize = 100
	} else if pageSize < 1 {
		pageSize = 1
	}

	logger.Infof("QueryTaskAPI: %v %v %v %v\n", UserID, taskType, pageNO, pageSize)

	tasks, err := services.QueryTaskService(UserID, taskType, uint(pageNO), uint(pageSize))
	if err != nil {
		gocommon.HttpErr(w, http.StatusOK, -1, err.Error())
		logger.Errorf("QueryTaskAPI ERR: %v %v %v\n", UserID, taskType, err.Error())
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, tasks)
	logger.Infof("QueryTaskAPI OK: %v %v %v\n", UserID, taskType)

	return
}
