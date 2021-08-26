package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

// version and tag of the code
var gitVersion, gitTag string

// Info function returns version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("PodManager git=%s tag=%s go=%s date=%s", gitVersion, gitTag, goVersion, tstamp)
}

func main() {
	var config string
	flag.StringVar(&config, "config", "", "config file name")
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)

	}
	server(config)
}
