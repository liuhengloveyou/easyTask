package common

import (
	"encoding/json"
	"fmt"
	"strconv"

	gocommon "github.com/liuhengloveyou/go-common"
)

type proxyResp struct {
	Code   int    `json:"code"`
	Ports  []int  `json:"port"`
	Domain string `json:"domain"`
	User   string `json:"authuser"`
	Pass   string `json:"authpass"`
}

func GetProxyURL() (domain, port, auth string, err error) {
	_, body, err := gocommon.GetRequest("http://14.18.242.38:8095/open?api=pkajuzne&area=440100", map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	})

	if err != nil {
		Logger.Sugar().Errorf("GetProxyURL ERR: ", err)
		return "", "", "", err
	}

	var rst proxyResp
	if err = json.Unmarshal(body, &rst); err != nil {
		Logger.Sugar().Errorf("GetProxyURL response ERR: ", string(body))
		return "", "", "", err
	}

	if rst.Code != 200 {
		Logger.Sugar().Errorf("GetProxyURL response.code ERR: ", string(body))
		return "", "", "", err
	}

	Logger.Sugar().Debugf("GetProxyURL: %v\n", string(body))

	return rst.Domain, strconv.Itoa(rst.Ports[0]), rst.User + ":" + rst.Pass, nil
}

func CloseProxyURL(port string) {
	_, body, err := gocommon.GetRequest(fmt.Sprintf("http://14.18.242.38:8095/close?api=pkajuzne&&port=%s", port), map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	})

	if err != nil {
		Logger.Sugar().Errorf("CloseProxyURL ERR: %v", err)
		return
	}

	Logger.Sugar().Debugf("CloseProxyURL: %s\n", string(body))

	return
}
