package rappers

import (
	"fmt"
	"strings"
)

// 任务描述信息
type DaikanTaskInfo struct {
	Rid        string `json:"rid"`     // 记录ID
	WebSite    string `json:"website"` // 网课网站名
	School     string `json:"school"`  // 学校名
	Student    string `json:"student"` // 学生学号
	StudentPWD string `json:"pwd"`     // 登录网课网站密码
	Course     string `json:"course"`  // 课程标题
}

func (p *DaikanTaskInfo) GetRid() string {
	return p.Rid
}

// 学习通/超星 湖南涉外 1841410112 3993831 asdf
func (p *DaikanTaskInfo) FromString(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) != 5 {
		return fmt.Errorf("format ERR")
	}

	p.WebSite = fields[0]
	p.School = fields[1]
	p.Student = fields[2]
	p.StudentPWD = fields[3]
	p.Course = fields[4]

	p.Rid = p.WebSite + "~" + p.School + "~" + p.Student + "~" + p.Course

	return nil
}

type DaikanRapper struct {
}

func NewDaikanRapper() Rapper {
	return &DaikanRapper{}
}

func (p *DaikanRapper) NewTaskInfo() TaskInfoI {
	return &DaikanTaskInfo{}
}

func (this *DaikanRapper) Run() {
	return
}
