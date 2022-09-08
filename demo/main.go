package main

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"time"
)

func main() {
	configure.SetDefault("FSS.Server", "http://127.0.0.1:888")
	fs.Initialize[StartupModule]("fss client")
	//initData()
	for {
		time.Sleep(time.Hour)
	}
}
