package main

import (
	app "kblswitch/internal"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)

	myApp := app.NewApp()
	myAppStop := myApp.Start()
	defer func() {
		if err := recover(); err != nil {
			myAppStop(err.(error))
			os.Exit(1)
		}
	}()
	myApp.KeepAlive()
	myAppStop(nil)
}
