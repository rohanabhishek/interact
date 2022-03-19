// This file is main launch for the server
package main

import (
	"flag"
	"github.com/golang/glog"
	web "interact/server/services/web"
	"sync"
)

var (
	addr = flag.String("addr", ":8080", "Server runs on this Address")
	// waitGroup to wait without ending the server
	waitGroup sync.WaitGroup
)

func main() {
	// Print logs into the files and stdouterr
	flag.Lookup("logtostderr").Value.Set("true")
	// Logs to this directory
	// flag.Lookup("log_dir").Value.Set("logs/")
	flag.Parse()

	webServer := web.NewWebServer(addr)
	webServer.Handlers()

	go webServer.Run()
	defer glog.Error("Server ended")

	waitGroup.Add(1)
	waitGroup.Wait()
}
