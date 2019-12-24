package api

import (
	"context"
	"net/http"
	"time"

	"github.com/liuhengloveyou/easyTask/common"

	gocommon "github.com/liuhengloveyou/go-common"
	passport "github.com/liuhengloveyou/passport/face"
	"github.com/liuhengloveyou/passport/sessions"
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func InitHttpApi(addr string) (handler *TaskHttpServer) {
	logger = common.Logger.Sugar()
	handler = &TaskHttpServer{}

	if addr != "" {
		//http.Handle("/monitor", &MonitorHandler{})
		http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./static/"))))
		http.Handle("/api", handler)
		s := &http.Server{
			Addr:           addr,
			ReadTimeout:    10 * time.Minute,
			WriteTimeout:   10 * time.Minute,
			MaxHeaderBytes: 1 << 20,
		}
		if err := s.ListenAndServe(); err != nil {
			panic(err)
		}
	}

	return
}

type TaskHttpServer struct {
	AuthFilter func(w http.ResponseWriter, r *http.Request) (*sessions.Session, bool)
}

func (p *TaskHttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 跨域资源共享
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)
	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}

	if p.AuthFilter != nil {
		sess, auth := p.AuthFilter(w, r)
		if false == auth {
			return // 没有登录
		}
		if sess == nil {
			return
		}
		r = r.WithContext(context.WithValue(context.Background(), "session", sess))
	}

	api := r.Header.Get("X-API")
	logger.Debugf("task api: %v\n", api)
	switch api {
	case "/task/add":
		AddTaskAPI(w, r) // 添加一个任务
	case "/task/add/batch":
		AddTaskBatchAPI(w, r) // 批量添加任务
	case "/task/delete":
		DeleteTaskAPI(w, r) // 删除自己的任务
	case "/task/query":
		QueryTaskAPI(w, r) // 查询任务详情
	case "/task/distribute":
		GetTaskAPI(w, r) // 分发任务
	case "/task/update":
		UpdateTaskAPI(w, r) // 更新任务
	case "/task/redo":
		RedoTaskAPI(w, r) // 重做
	default:
		gocommon.HttpErr(w, http.StatusNotFound, 0, "")
		return
	}

	return
}

func AuthFilter(w http.ResponseWriter, r *http.Request) (sess *sessions.Session, auth bool) {
	sess, auth = passport.AuthFilter(w, r)
	logger.Debug("session:", sess, auth)

	if auth == false && sess == nil {
		gocommon.HttpErr(w, http.StatusUnauthorized, -1, "请登录")
		return
	} else if auth == false && sess != nil {
		gocommon.HttpErr(w, http.StatusUnauthorized, -1, "您没有权限")
		return
	}

	return
}
