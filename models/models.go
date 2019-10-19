package models

import (
	"database/sql"

	"github.com/liuhengloveyou/easyTask/common"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

const (
	TaskStatusNew = iota
	TaskStatusSend
	TaskStatusOK
	TaskStatusERR

	TaskStatusEND
)

var (
	logger *zap.SugaredLogger
)

type DAOInterface interface {
	Insert(tx *sql.Tx) (sql.Result, error)
	Update(tx *sql.Tx) (sql.Result, error)
	Delete(tx *sql.Tx) (sql.Result, error)
}

func init() {
	logger = common.Logger.Sugar()
}

func Insert(tx *sql.Tx, model DAOInterface) (rst sql.Result, e error) {
	var _tx *sql.Tx

	if tx != nil {
		_tx = tx
	} else {
		if _tx, e = common.DB.Begin(); e != nil {
			return nil, e
		}

		defer Rollback(_tx)
	}

	rst, e = model.Insert(tx)

	if tx != _tx {
		if e = _tx.Commit(); e != nil {
			logger.Errorf("tx.Commit ERR: ", e.Error())
		}
	}

	return
}

func Delete(tx *sql.Tx, model DAOInterface) (rst sql.Result, e error) {
	var _tx *sql.Tx

	if tx != nil {
		_tx = tx
	} else {
		if _tx, e = common.DB.Begin(); e != nil {
			return nil, e
		}

		defer Rollback(_tx)
	}

	rst, e = model.Delete(tx)

	if tx != _tx {
		if e = _tx.Commit(); e != nil {
			logger.Errorf("tx.Commit ERR: ", e.Error())
		}
	}

	return
}

func Update(tx *sql.Tx, model DAOInterface) (rst sql.Result, e error) {
	var _tx *sql.Tx

	if tx != nil {
		_tx = tx
	} else {
		if _tx, e = common.DB.Begin(); e != nil {
			return nil, e
		}

		defer Rollback(_tx)
	}

	rst, e = model.Update(_tx)

	if tx != _tx {
		if e = _tx.Commit(); e != nil {
			logger.Errorf("tx.Commit ERR: ", e.Error())
		}
	}

	return
}

func BeginTx() (*sql.Tx, error) {
	return common.DB.Begin()
}

// defer Rollback(_tx)
func Rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		logger.Errorf("tx.Rollback ERR: ", err.Error())
	}
}
