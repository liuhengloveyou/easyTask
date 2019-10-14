package models

import (
	"database/sql"
	"fmt"

	gocommon "github.com/liuhengloveyou/go-common"
)

type Task struct {
	ID         int64         `json:"id"`
	Tid        string        `json:"tid"`                          // 任务ID
	Rid        string        `json:"rid"`                          // 记录ID
	TaskType   string        `json:"task_type"`                    // 任务类型
	TaskInfo   gocommon.JSON `json:"task_info"`                    // 任务内容
	Stat       int           `json:"stat"`                         // 任务状态
	AddTime    sql.NullTime  `json:"add_time" db:"add_time"`       // 添加时间
	UpdateTime sql.NullTime  `json:"update_time" db:"update_time"` // 更新时间
	Rapper     string        `json:"rapper"`
	Client     string        `json:"client"`
	Remark     string        `json:"remark"`
}

func (p *Task) Insert() (id int64, e error) {
	rdb := db.Create(p)

	return p.ID, rdb.Error
}

func (p *Task) Query(num int) (tasks []Task, e error) {
	rdb := db.Where("task_type = ?", p.TaskType).Limit(num).Find(&tasks)

	return tasks, rdb.Error
}

func (p *Task) Update() (e error) {
	db.Model(p).Where("active = ?", true).Update("name", "hello")

	return nil
}
