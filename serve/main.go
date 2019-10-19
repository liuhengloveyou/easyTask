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

	"github.com/liuhengloveyou/easyTask/api"
	"github.com/liuhengloveyou/easyTask/common"
	gocommon "github.com/liuhengloveyou/go-common"
)

var Sig string

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
	sigHandler()
}

func main() {
	// go RapperCleaner()

	fmt.Println("easytask GO...", common.ServeConfig.Addr)
	panic(api.InitHttpApi(common.ServeConfig.Addr))
}
