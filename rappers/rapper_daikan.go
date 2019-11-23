package rappers

import (
	"fmt"
	"strings"
)

// 任务描述信息
type DaikanTaskInfo struct {
	Rid        string `json:"rid"`     // 记录ID
	WebSite    string `json:"website"` // 网课网站名
	Account     string `json:"account"`  // 账号
	PWD string `json:"pwd"`     // 登录网课网站密码
	Course     string `json:"course"`  // 课程标题
}

func (p *DaikanTaskInfo) GetRid() string {
	return p.Rid
}

// 学习通/超星 湖南涉外 1841410112 3993831 asdf
func (p *DaikanTaskInfo) FromString(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) != 4 {
		return fmt.Errorf("数据格式错误")
	}

	p.WebSite = fields[0]
	p.Account = fields[1]
	p.PWD = fields[2]
	p.Course = fields[3]

	p.Rid = p.WebSite + "~" + p.Account + "~" + p.Course

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
