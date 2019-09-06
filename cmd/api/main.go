package main

import (
	container_ai "bucket_file/cmd/api/container"
	"flag"
	"log"
	"runtime"
)

var (
	container = container_ai.NewContainer()
	config    string
)

func main() {
	flag.Parse()
	err := container.Setup(config)
	if err != nil {
		log.Fatal(err)
	}

	CreateFile()
	createAdmin()
	Service()
	
}

func init() {
	flag.StringVar(&config, "config", "config.toml", "path of file config")
	runtime.GOMAXPROCS(runtime.NumCPU())
}
