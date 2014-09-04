package common

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBmysql struct {
	Url string
	conn       *sql.DB
}

func (this *DBmysql) Init(url string) (*DBmysql, error) {
	err := fmt.Errorf("Init mysql connection witch nil url.")
	if url == "" {
		return nil, err
	}

	if this.conn, err = sql.Open("mysql", url); err != nil {
		return nil, err
	}

	return this, nil
}

func (this *DBmysql) Query(sqlStr string, args ...interface{}) (rst []map[string]string, err error) {
	var (
		stmt *sql.Stmt = nil
		rows *sql.Rows = nil
	)
	
	stmt, err = this.conn.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(args...)
	if err != nil {
		return
	}
	defer rows.Close()
	
	var cols []string
	cols, err = rows.Columns()
	if err != nil {
		return
	}

	cvals := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(cols))
	for i := range cvals {
		scanArgs[i] = &cvals[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}

		tmap := make(map[string]string, len(cols))
		for i, col := range cvals {
			if col == nil {
				tmap[cols[i]] = ""
			} else {
				tmap[cols[i]] = string(col)
			}
		}
		rst = append(rst, tmap)
	}

	return
}

func (this *DBmysql) Insert(sqlStr string, args ...interface{}) (int64, error) {
	stmt, err := this.conn.Prepare(sqlStr)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}

	return rst.LastInsertId()
}

func (this *DBmysql) Update(sqlStr string, args ...interface{}) (int64, error) {
	stmt, err := this.conn.Prepare(sqlStr)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}

	return rst.RowsAffected()
}
