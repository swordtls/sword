package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sword/core"
)

var runModeMap = map[string]int{"s":1,"server":2,"c":3,"client":4,"ps":5,"proxy_server":6,"b":7,"base64p12":8}

func main() {
	runMode := flag.String("mode", "unknown", "run mode, s(server),c(client),ps(proxy_server),(b)base64p12")
	config := flag.String("conf", "./server.json", "input config file")
	daemon := flag.Bool("daemon", false, "run background")
	p12 := flag.String("p12","client.p12","p12 file name ")
	if _,ok := runModeMap[*runMode];ok && *daemon == true {
		if os.PathSeparator == 47 {
			if os.Getppid() != 1 {
				filePath, _ := filepath.Abs(os.Args[0])
				cmd := exec.Command(filePath, os.Args[1:]...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Start()
				return
			}
		}
	}
	flag.Parse()
	switch *runMode {
	case "s","server":
		core.Sever(*config)
	case "c","client":
		core.Client(*config)
	case "ps","proxy_server":
		core.ProxyServer(*config)
	case "b","base64p12":
		core.Base64P12(*p12)
	default:
		fmt.Println("Unknown run mode: sword -h")

	}
}

