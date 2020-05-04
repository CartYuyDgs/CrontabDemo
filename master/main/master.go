package main

import (
	"CrontabDemo/master"
	"flag"
	"fmt"
	"runtime"
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

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	return

ERR:
	fmt.Println(err)

}
