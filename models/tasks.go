package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/liuhengloveyou/easyTask/common"

	"github.com/didi/gendry/builder"
	gocommon "github.com/liuhengloveyou/go-common"
)

type Task struct {
	ID         int64         `json:"id" db:"id"`               // 任务ID
	UID        int64         `json:"uid" db:"uid"`             // 用户ID
	Rid        string        `json:"rid" db:"rid"`             // 记录ID
	TaskInfo   gocommon.JSON `json:"task_info" db:"task_info"` // 任务内容
	TaskType   string        `json:"task_type" db:"task_type"` // 任务类型
	Stat       string        `json:"stat" db:"stat"`           // 任务状态
	Rapper     string        `json:"name" db:"rapper"`
	AddTime    *time.Time    `json:"addTime" db:"add_time"`       // 添加时间
	UpdateTime *time.Time    `json:"updateTime" db:"update_time"` // 更新时间
	Remark     string        `json:"remark" db:"remark"`
}

func (p *Task) Insert() (id int64, e error) {
	var rst sql.Result

	table := "tasks"
	data := map[string]interface{}{
		"rid":         p.Rid,
		"task_type":   p.TaskType,
		"task_info":   p.TaskInfo,
		"stat":        p.Stat,
		"add_time":    time.Now(),
		"update_time": sql.NullTime{Valid: true, Time: time.Now()},
	}

	if p.UID > 0 {
		data["uid"] = p.UID
	}

	sqlStr, valArr, err := builder.BuildInsert(table, []map[string]interface{}{data})
	if err != nil {
		return -1, err
	}

	logger.Debug("Task.Insert sql: ", sqlStr, valArr, err)

	rst, e = common.DB.Exec(sqlStr, valArr...)
	if e == nil {
		id, e = rst.LastInsertId()
	}

	return
}

func (p *Task) Query(pageNO, pageSize uint) (tasks []Task, e error) {
	table := "tasks"
	selectFields := []string{"id", "uid", "rid", "task_type", "task_info", "stat", "update_time"}
	where := map[string]interface{}{
		"_orderby": "id asc",
		"_limit":   []uint{(pageNO - 1) * pageSize, pageSize},
	}

	if p.TaskType != "" {
		where["task_type"] = p.TaskType
	}
	if p.Stat != "" {
		where["stat"] = p.Stat
	}
	if p.UID > 0 {
		where["uid"] = p.UID
	}

	cond, vals, err := builder.BuildSelect(table, where, selectFields)
	logger.Info("Task.Select sql: ", cond, vals, err)

	if e = common.DB.Select(&tasks, cond, vals...); e != nil {
		return
	}

	return tasks, nil
}

func (p *Task) Update() (row int64, e error) {
	var rst sql.Result

	table := "tasks"
	where := map[string]interface{}{
		"id": p.ID,
	}
	if p.UpdateTime != nil {
		where["update_time"] = p.UpdateTime

	}
	if p.UID > 0 {

	}

	update := make(map[string]interface{})
	if p.Stat != "" {
		update["stat"] = p.Stat
	}
	if p.Rapper != "" {
		update["rapper"] = p.Rapper
	}
	if p.Remark != "" {
		update["remark"] = p.Remark
	}

	sqlStr, valArr, err := builder.BuildUpdate(table, where, update)
	common.Logger.Sugar().Debug("Task.Update sql: ", sqlStr, valArr, err)
	if err != nil {
		common.Logger.Sugar().Errorf("Task.update BuildUpdate ERR: ", err.Error())
		return 0, fmt.Errorf("服务错误")
	}

	logger.Info("Task.update sql: ", sqlStr, valArr)
	rst, e = common.DB.Exec(sqlStr, valArr...)
	if e == nil {
		return rst.RowsAffected()
	}

	return
}
