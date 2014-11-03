package common

import (
	"os"
	"encoding/json"
)

var Sig string
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


