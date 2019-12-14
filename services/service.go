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

		sql := fmt.Sprintf("update tasks set stat='重做' where update_time < '%v' and stat != '完成' and stat != '新增' and stat != '出错'", time.Now().Add(30*time.Minute).Format("2006-01-02 15:04:05"))
		rst, e := common.DB.Exec(sql)
		if e == nil {
			row, _ = rst.RowsAffected()
		}
		logger.Infof("timeOutTask: %v %v\n", row, e)

		time.Sleep(time.Minute)
	}
}
