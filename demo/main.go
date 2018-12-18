package main

import (
	"os"

	"github.com/zealllot/configor"
)

func main() {
	here, _ := os.Getwd()
	configPath := here + "/test"

	configor.Load(configPath)
	select {}
}
