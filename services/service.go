package services

import (
	"fmt"
	"time"

	"github.com/liuhengloveyou/easyTask/common"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	logger = common.Logger.Sugar()

	go timeOutTask()
}

func timeOutTask() {
	for {
		var row int64 = 0

		sql := fmt.Sprintf("update tasks set stat='末处理' where update_time < '%v' and stat != '完成' and stat != '末处理' and stat != '失败'", time.Now().Add(-1*time.Hour).Format("2006-01-02 15:04:05"))
		rst, e := common.DB.Exec(sql)
		if e == nil {
			row, _ = rst.RowsAffected()
		}
		logger.Infof("timeOutTask: %v %v\n", row, e)

		time.Sleep(time.Minute)
	}
}
