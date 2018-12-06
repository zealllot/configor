package main

import (
	"os"

	"github.com/zealllot/configer"
)

func main() {
	here, _ := os.Getwd()
	configPath := here + "/test"

	//fmt.Println(here)
	configer.AddConfigPath(configPath)
	select {}
}
