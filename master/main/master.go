package main

import (
	"CrontabDemo/master"
	"fmt"
	"runtime"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	initEnv()

	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	return

ERR:
	fmt.Println(err)

}
