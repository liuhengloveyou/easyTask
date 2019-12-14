package rappers

import (
	"fmt"
	"strings"
)

// 任务描述信息
type DaikanTaskInfo struct {
	Rid     string `json:"rid"`     // 记录ID
	WebSite string `json:"website"` // 网课网站名
	School  string `json:"school"`  // 学校名
	Account string `json:"account"` // 账号/学号
	PWD     string `json:"pwd"`     // 登录网课网站密码
	Course  string `json:"course"`  // 课程标题
}

func (p *DaikanTaskInfo) GetRid() string {
	return p.Rid
}

// 智慧树 广东司法警官职业学院 18578744736 123456 课程标题
// 智慧树 18289264379 123456 课程标题
// 学习通/超星 湖南涉外 1841410112 3993831 asdf
func (p *DaikanTaskInfo) FromString(raw string) error {
	fields := strings.Fields(raw)
	if len(fields) == 4 {
		p.WebSite = fields[0]
		p.Account = fields[1]
		p.PWD = fields[2]
		p.Course = fields[3]
		p.Rid = p.WebSite + "~" + p.Account + "~" + p.Course
	} else if len(fields) == 5 {
		p.WebSite = fields[0]
		p.School = fields[1]
		p.Account = fields[2]
		p.PWD = fields[3]
		p.Course = fields[4]
		p.Rid = p.WebSite + "~" + p.Account + "~" + p.Course
	} else {
		return fmt.Errorf("数据格式错误")
	}

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
