package models

import (
	"database/sql"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/didi/gendry/builder"
)

type Task struct {
	ID         int64         `json:"id" db:"id"`                   // 任务ID
	Rid        string        `json:"rid" db:"rid"`                 // 记录ID
	TaskType   string        `json:"task_type" db:"task_type"`     // 任务类型
	TaskInfo   gocommon.JSON `json:"task_info" db:"task_info"`     // 任务内容
	Stat       int           `json:"stat" db:"stat"`               // 任务状态
	AddTime    time.Time     `json:"add_time" db:"add_time"`       // 添加时间
	UpdateTime time.Time     `json:"update_time" db:"update_time"` // 更新时间
	Rapper     string        `json:"rapper" db:"rapper"`
	Client     string        `json:"client" db:"client"`
	Remark     string        `json:"remark" db:"remark"`
}

func (p *Task) Insert() (id int64, e error) {
	var rst sql.Result

	table := "tasks"
	data := map[string]interface{}{
		"rid":       p.Rid,
		"task_type": p.TaskType,
		"task_info": p.TaskInfo,
		"stat":      TaskStatusNew,
		"add_time":  time.Now(),
	}

	sqlStr, valArr, err := builder.BuildInsert(table, []map[string]interface{}{data})
	if err != nil {
		return -1, err
	}

	logger.Debug("Task.Insert sql: ", sqlStr, valArr, err)

	rst, e = db.Exec(sqlStr, valArr...)
	if e == nil {
		id, e = rst.LastInsertId()
	}

	return
}

func (p *Task) Query(num int) (tasks []Task, e error) {
	table := "tasks"
	selectFields := []string{"id", "rid", "task_type", "task_info", "stat", "update_time"}
	where := map[string]interface{}{
		"task_type": p.TaskType,
		"stat": 0,
		"_orderby": "id asc",
		"_limit":   []uint{0, uint(num)},
	}

	cond, vals, err := builder.BuildSelect(table, where, selectFields)
	logger.Info("Task.Select sql: ", cond, vals, err)

	if e = db.Select(&tasks, cond, vals...); e != nil {
		return
	}

	return tasks, nil
}

func (p *Task) Update() (row int64, e error) {
	var rst sql.Result

	sqlStr := "UPDATE tasks SET stat=?, remark=? WHERE (id=? AND update_time=?)"
	where := []interface{}{p.Stat, p.Remark, p.ID, p.UpdateTime}

	logger.Info("Task.update sql: ", sqlStr, where)
	rst, e = db.Exec(sqlStr, where...)
	if e == nil {
		return rst.RowsAffected()
	}

	return
}