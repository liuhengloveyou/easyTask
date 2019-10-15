package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/liuhengloveyou/easyTask/common"
	"github.com/liuhengloveyou/easyTask/rappers"

	gocommon "github.com/liuhengloveyou/go-common"
)

type taskInfo struct {
	Tid  string // 任务ID
	Rid  string // 记录ID
	Info string // 任务内容
}

var (
	Sig string
)

func sigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		s := <-c
		Sig = "service is suspend ..."
		fmt.Println("Got signal:", s)
	}()
}

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	gocommon.SingleInstane(common.ServeConfig.PID) // 单例

	rand.Seed(time.Now().UTC().UnixNano())

	if err := common.InitClient(); err != nil {
		panic(err)
	}

	sigHandler()
}

func main() {

	// 向服务器打招乎
	//if err := sayHiToServ(); err != nil {
	//	panic(err)
	//}

	for i := 0; i < common.ClientConfig.Flow; i++ {
		rapper, err := rappers.NewRapper(common.ClientConfig.TaskType)
		if err != nil {
			panic(err)
		}

		go rapper.Run()
	}

	for {
		time.Sleep(5 * time.Second)
	}
}
