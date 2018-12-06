package main

import (
	//"github.com/zealllot/configer"
	"os"
	"fmt"
	"path"
)

func main() {
	here,_:=os.Getwd()
	configPath:=here+"/demo/test"

	fd,err:=os.Stat(configPath)
	if err != nil {
		panic(err)
	}
	if fd.IsDir(){
		fmt.Println("大骗子")
		return
	}

	fmt.Println("父目录:",path.Dir(configPath))
	//fmt.Println(here)
	//configer.AddConfigDirectory(here+"/test")
	//select {}
}
