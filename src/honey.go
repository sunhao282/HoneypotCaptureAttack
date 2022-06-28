package main

import (
	"honey/src/setting"
	"os"
	"runtime"
)

func main() {
	//设置最大CPU数
	runtime.GOMAXPROCS(runtime.NumCPU())
	args := os.Args
	if args == nil || len(args) < 2 {
		setting.Help()
	} else {
		if args[1] == "-h" {
			setting.Help()
			//} else if args[1] == "init" {
			//	setting.Init()
		} else if args[1] == "-v" {
			setting.Version()
		} else if args[1] == "-r" {
			setting.Run()
		} else {
			setting.Help()
		}
	}
	select {}
}
