package api

import (
	"net/http"
	"time"

	"github.com/liuhengloveyou/easyTask/common"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
	router *httprouter.Router
)

func InitHttpApi(addr string) error {
	router = httprouter.New()
	logger = common.Logger.Sugar()

	router.PUT("/addtask/:type/:count", AddTaskAPI)
	router.GET("/querytask", QueryTaskAPI)
	router.POST("/updatetask", UpdateTaskAPI)

	//http.Handle("/sayhi", &SayhiHandler{})
	//http.HandleFunc("/beat", HandleBeat)
	//
	//http.Handle("/monitor", &MonitorHandler{})
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./static/"))))

	// 用户
	//router.POST("/api/user", UpdateUserInfo)      // 更新用户信息
	//router.GET("/api/user", GetMyInfo)            // 查询用户个人信息
	//router.GET("/api/user/open", GetUserInfoOpen) // 查询用户公开信息

	// passport
	// http.Handle("/user", &passport.HttpServer{})

	// root
	http.Handle("/", &Server{})

	s := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

type Server struct{}

func (p *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer logger.Sync()

	// 跨域资源共享
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8100")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "X-API, X-REQUEST-ID, X-API-TRANSACTION, X-API-TRANSACTION-TIMEOUT, X-RANGE, Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	//URL, err := url.ParseRequestURI(r.RequestURI)
	//if err != nil {
	//	logger.Error("url ERR: ", err)
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}

	//sess, auth := AuthFilter(w, r)
	//
	//if false == auth {
	//	return // 没有登录
	//}
	//
	//if sess != nil {
	//	r = r.WithContext(context.WithValue(context.Background(), "session", sess))
	//}

	router.ServeHTTP(w, r)

	return
}

//
//func AuthFilter(w http.ResponseWriter, r *http.Request) (sess *sessions.Session, auth bool) {
//	sess, auth = passport.AuthFilter(w, r)
//	logger.Debug("session:", sess, auth)
//
//	if auth == false && sess == nil {
//		gocommon.HttpErr(w, http.StatusOK, -1, "请登录")
//		return
//	} else if auth == false && sess != nil {
//		gocommon.HttpErr(w, http.StatusOK, -1, "您没有权限")
//		return
//	}
//
//	return
//}
