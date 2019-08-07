package main

import (
	"Coot/utils/setting"
	"fmt"
	"os"
)

func main() {
	setting.RunWeb("localhost:9000")
	args := os.Args
	if args == nil || len(args) < 2 {
		setting.Help()
	} else {
		if args[1] == "help" || args[1] == "--help" {
			setting.Help()
		} else if args[1] == "init" || args[1] == "--init" {
			setting.Init()
		} else if args[1] == "version" || args[1] == "--version" {
			fmt.Println("0.1")
		} else if args[1] == "run" || args[1] == "--run" {
			if len(args) >= 3 {
				setting.RunWeb(args[2])
			} else {
				setting.RunWeb("localhost:9000")
			}
		} else {
			setting.Help()
		}
	}
}
