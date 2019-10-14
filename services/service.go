package services

import (
	"github.com/liuhengloveyou/easyTask/common"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func init () {
	logger = common.Logger.Sugar()
}