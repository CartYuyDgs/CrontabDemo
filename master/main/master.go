package main

import (
	"CrontabDemo/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	configFile string
)

func initArgs() {
	//master-config ./master.json
	flag.StringVar(&configFile, "config", "./master.json", "指定配置文件")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	initEnv()
	initArgs()

	if err = master.InitConfig(configFile); err != nil {
		goto ERR
	}

	if err = master.InitLogManger(); err != nil {
		goto ERR
	}

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}
	return

ERR:
	fmt.Println(err)

}
