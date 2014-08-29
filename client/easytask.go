package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/golang/glog"
)

type rapper interface {
	run() // 开始任务
}

type rapperType func() rapper

type taskInfo struct {
	Tid  string // 任务ID
	Rid  string // 记录ID
	Info string // 任务内容
}

var (
	confJson map[string]interface{}
	rappers  = make(map[string]rapperType)
	tasks    chan *taskInfo
)

func init() {
	runtime.GOMAXPROCS(8)

	r, err := os.Open("./app.conf")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&confJson)
	if err != nil {
		panic(err)
	}

	tasks = make(chan *taskInfo, int(confJson["flow"].(float64))*2)
}

func download(url, fn string) error {
	if url == "" || fn == "" {
		return fmt.Errorf("param err")
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errRst error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errRst = fmt.Errorf("%d\r\n%s", resp.StatusCode, err.Error())
		} else {
			errRst = fmt.Errorf("%d\r\n%s", resp.StatusCode, string(body))
		}
		return errRst
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func upload(url, fn string, para *map[string]string) ([]byte, error) {
	if url == "" || fn == "" {
		return nil, fmt.Errorf("param err")
	}

	fp, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if para != nil {
		for k, v := range *para {
			writer.WriteField(k, v)
		}
	}

	_, err = writer.CreateFormFile("file", fn)
	if err != nil {
		return nil, err
	}

	boundary := writer.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	reqReader := io.MultiReader(body, fp, closeBuf)
	req, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.ContentLength = fi.Size() + int64(body.Len()) + int64(closeBuf.Len())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		var errRst error
		if err != nil {
			errRst = fmt.Errorf("%v\r\n%s", resp.StatusCode, err.Error())
		} else {
			errRst = fmt.Errorf("%v\r\n%s", resp.StatusCode, respBody)
		}
		return nil, errRst
	}

	return respBody, nil
}

func getRequest(url string, para *map[string]string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("get URL nil")
	}

	if para != nil {
		url += "?_=_"
		for k, v := range *para {
			url += "&" + k + "=" + v
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d\r\n%s", resp.StatusCode, string(body))
	}

	return body, nil
}

var once sync.Once

func oneTask() *taskInfo {
	once.Do(func() {
		go func() {
			count := 0
			urlStr := fmt.Sprintf("%s/getask?name=%s&type=%s&num=%d", confJson["taskServ"], confJson["rappername"], confJson["tasktype"], int64(confJson["flow"].(float64))*2+3)
			for {
				resp, err := http.Get(urlStr)
				if err != nil {
					glog.Errorln(err, urlStr)
					time.Sleep(1 * time.Second)
					continue
				}

				body, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil || resp.StatusCode != 200 {
					glog.Errorln(err, resp.StatusCode, string(body))
					time.Sleep(1 * time.Second)
					continue
				}

				var taskJson []taskInfo
				err = json.Unmarshal(body, &taskJson)
				if err != nil {
					glog.Errorln(err, string(body))
					time.Sleep(1 * time.Second)
					continue
				}

				if 1 > len(taskJson) {
					glog.Errorln("gettask nil")
					time.Sleep(1 * time.Second)
					continue
				}

				for i := 0; i < len(taskJson); i++ {
					if taskJson[i].Tid == "" || taskJson[i].Rid == "" || taskJson[i].Info == "" {
						glog.Errorln("taskERR: ", taskJson[i])
						continue
					}
					count++
					glog.Infoln(count, taskJson[i])
					tasks <- &taskJson[i]
				}
			}
		}()
	})

	return <-tasks
}

func sayHiToServ() error {
	para := &map[string]string{"type": confJson["tasktype"].(string), "name": confJson["rappername"].(string)}
	resp, err := getRequest(confJson["taskServ"].(string)+"/sayhi", para)
	if err != nil {
		return err
	}

	if string(resp) != "OK" {
		return fmt.Errorf("%s", resp)
	}

	return nil
}

func sendBeat() error {
	urlStr := fmt.Sprintf("%s/beat?type=%s&name=%s", confJson["taskServ"].(string), confJson["tasktype"].(string), confJson["rappername"].(string))
	_, err := getRequest(urlStr, nil)
	if err != nil {
		return err
	}

	return nil
}

func register(name string, one rapperType) {
	if one == nil {
		panic("register rapper nil")
	}
	if _, dup := rappers[name]; dup {
		panic("register duplicate for " + name)
	}
	rappers[name] = one
}

func NewRapper(typeName string) (rapper, error) {
	newFun, ok := rappers[typeName]
	if ok != true {
		return nil, fmt.Errorf("no rapper types " + typeName)
	}

	return newFun(), nil
}

func main() {
	flag.Parse()

	// 向服务器打招乎
	if err := sayHiToServ(); err != nil {
		panic(err)
	}

	for i := 0; i < int(confJson["flow"].(float64)); i++ {
		one, err := NewRapper(confJson["tasktype"].(string))
		if err != nil {
			panic(err)
		}
		go one.run()
	}

	for {
		time.Sleep(5 * time.Second)
		sendBeat()
	}
}
