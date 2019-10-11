package common

import (
	"encoding/json"
	"os"
)

var ConfJson map[string]interface{} // 系统配置信息

func init() {
	r, err := os.Open("./app.conf")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&ConfJson); err != nil {
		panic(err)
	}
}
