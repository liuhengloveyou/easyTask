package models

import (
	"github.com/liuhengloveyou/easyTask/common"

	"go.uber.org/zap"
	"github.com/jinzhu/gorm"
)


var (
	logger *zap.SugaredLogger
	db     *gorm.DB
)

func InitDB() error {
	logger = common.Logger.Sugar()

	var err error
	if db, err = gorm.Open("mysql", common.ServeConfig.Mysql); err != nil {
		return err
	}

	db.SetLogger(&common.MyLogger{})
	db.LogMode(true)

	return nil
}