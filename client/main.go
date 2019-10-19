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
	rand.Seed(time.Now().UTC().UnixNano())

	if err := common.InitClient(); err != nil {
		panic(err)
	}

	gocommon.SingleInstane(common.ClientConfig.PID) // 单例

	sigHandler()
}

func main() {
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
