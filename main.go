package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"os/exec"
	"io/ioutil"
	"strings"

	//"github.com/shirou/gopsutil/process"
	"gopkg.in/yaml.v1"
)

type config struct {
	UsagiProspector []Prospector `yaml:"usagiprospector"`
}

type Prospector struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Param string `yaml:"param"`
}

func checkBinaryExists(b string) bool {
	bin, err := exec.LookPath(b)
	if err != nil {
		panic(err)
	}
	if b != "" && len(string(b)) > 0 {
		fmt.Println(fmt.Sprintf("%s binary located at %s\n\n", b, bin))
		return true
	} else {
		fmt.Println(fmt.Sprintf("%s binary is not found\n\n", b))
		return false
	}
}

func readConfig() config {
	var c config
	f, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println("error opening config file")
	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		fmt.Println("error unmarshalling config file", err)
	}
	return c
}

func checkIfUp(n string) bool {
	cmd := strings.Replace("ps ux | awk '/PROCESS/ && !/awk/ {print $2}'", "PROCESS", strings.Replace(n, "/", `\/`, -1), -1)
	fmt.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
			fmt.Println(err)
	}
	fmt.Println(out)
	if len(out) < 1 {
		fmt.Printf("%d", len(out))
		return false
	}

	return true
}

func start(n, args string) {
	binary, err := exec.LookPath(n)
	if err != nil {
		panic(err)
	}
	args += " &"
	exec := exec.Command("sh", "-c", binary, args).Run()
	if exec != nil {
		panic(exec)
	}
}

func main(){

	fmt.Println("Starting USAGI - Zombifying Daemon")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	shutdown := false

	cfg := readConfig()

	/* for _, v := range c.UsagiProspector {
		fmt.Printf("name: %s\n", v.Name)
		fmt.Printf("path: %s\n", v.Path)
		fmt.Printf("param: %s\n", v.Param)
	} */

	OUTTER:
	for {
		if shutdown {
			break
		}

		select {
		case s, ok := <- c:
			if ok {
				fmt.Println("\nUSAGI is shutting down...\n")
				fmt.Println(fmt.Sprintf("\n\nReceived signal: %x\n\n", s))
				//shutdown sequence started
				shutdown = true
				//close go channel for os signal
				c = nil
				continue OUTTER
			}
		default:
			for _, v := range cfg.UsagiProspector {
				run := checkIfUp(v.Path)
				if !run {
					//Start process since its not running
					go start(v.Path, v.Param)
				}
			}
			//finished checking, go to sleep
			time.Sleep(10000 * time.Millisecond)
		}
	}

	fmt.Println("USAGI has shutdown gracefully. Sayonara!")
}
