// This file is main launcher for the server
package main

import (
	"flag"
	web "interact/server/services/web"
	"sync"

	"github.com/golang/glog"
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

	webServer := web.NewWebServer(*addr)
	waitGroup.Add(1)

	go webServer.Run()
	defer glog.Error("Server ended")

	waitGroup.Wait()
}
