package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/DeanThompson/ginpprof"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/viper"
	"github.com/zhj0811/fabric-normal/apiserver/router"
	"github.com/zhj0811/fabric-normal/apiserver/sdk"
)

var configPath = flag.String("configPath", "./client_sdk.yaml", "config path")
var logger = flogging.MustGetLogger("service")

func main() {
	var err error

	runtime.GOMAXPROCS(runtime.NumCPU())
	// init sdk
	logger.Infof("Using config file %s", *configPath)
	err = sdk.InitSDKs(*configPath)
	if err != nil {
		logger.Errorf("Init sdk failed %s", err.Error())
		panic(err)
	}
	logger.Info("Init sdk success.")
	// init gin http server
	gin.SetMode(gin.ReleaseMode)
	r := router.GetRouter()
	ginpprof.Wrapper(r) // for debug
	listenPort := viper.GetInt("apiserver.listenport")
	if listenPort == 0 {
		listenPort = 8888
	}
	logger.Debug("The listen port is ", listenPort)
	server := endless.NewServer(fmt.Sprintf(":%d", listenPort), r)

	// save pid file
	server.BeforeBegin = func(add string) {
		pid := syscall.Getpid()
		logger.Criticalf("Actual pid is %d", pid)
		pidFile := "apiserver.pid"
		if checkFileIsExist(pidFile) {
			os.Remove(pidFile)
		}
		if err := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0666); err != nil {
			logger.Fatalf("Api server write pid file failed! err:%v", err)
		}
	}

	err = server.ListenAndServe()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			logger.Errorf("%v", err)
		} else {
			logger.Errorf("Api server start failed! err:%v", err)
			panic(err)
		}
		panic(err)
	}
	panic(err)
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
