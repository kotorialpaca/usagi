package main

import (
	"fmt"
	"log"
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
	Waittime time.Duration `yaml:"waittime"`
	UsagiProspector []Prospector `yaml:"usagiprospector"`
}

type Prospector struct {
	Name string `yaml:"name"`
	Search string `yaml:"search"`
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
	cmd := strings.Replace("ps -ef | awk '/PROCESS/ && !/awk/ {print $2}'", "PROCESS", strings.Replace(n, "/", `\/`, -1), -1)
	//check command w/ below
	//fmt.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
			log.Println(err)
			//fmt.Println(err)
	}
	//check PID with below, if exist
	//else process isnt running
	//fmt.Println(out)
	if len(out) < 1 {
		//checking length of output, if not running it wil yield 0
		//fmt.Printf("%d", len(out))
		log.Printf("Process %s is not running, reviving process\n", n)
		return false
	}
	
	log.Printf("Process %s is running, moving on to next task\n", n)

	return true
}

func start(n, args string) {
	binary, err := exec.LookPath(n)
	if err != nil {
		panic(err)
	}
	args += " &"
	log.Printf("Attempting to start %s\n", n)
	exec := exec.Command("sh", "-c", binary, args).Start()
	if exec != nil {
		panic(exec)
	}
}

func setDefault(c config) config {
	if c.Waittime == 0 {
		c.Waittime = 60000
	}
	return c
}

func printUsagi(){
	
	fmt.Println(`----------------------------------------------------------`)
	fmt.Println(`                g,                        g,`)
	fmt.Println(`              vQmpg                   _vgQp,`)
	fmt.Println(`              dQQQmp,                vgWQQQf`)
	fmt.Println(`             =mQQQQQms             _qWQQQQQ>`)
	fmt.Println(`              dQQQQQQms           _qWQQQQQQf`)
	fmt.Println(`              )QQQQQQQms         _qWQQQQQQE'`)
	fmt.Println(`               mQQQQQQQEggggggggg)QQQQQQQ@f`)
	fmt.Println(`               ]$QQQQQQQnnnnnnnnnmQQQQQQQf`)
	fmt.Println(`                ]$QQQQQQEnnnnnnnnQQQQQQQf`)
	fmt.Println(`                gnVQQQQQmnnnnnnndQQQQQ@vs,`)
	fmt.Println(`               %nnndQQQQQQQQQQQQQQQQQVvnnn,`)
	fmt.Println(`              jonnngWQQQQQQQQQQQQQQQQmpvnnn`)
	fmt.Println(`              oonqmQQ@WQQQQQQQQQQQQ@QQQmnnnL`)
	fmt.Println(`              onnQQQQQv3H$QQQQQQVVndQQQQEnnc`)
	fmt.Println(`              nndQQQQEQQnmQQQQQmEmQEQQQQQnn(`)
	fmt.Println(`              {nmQQQQQggQQQ@VVQQQmgQWQQQQnn'`)
	fmt.Println(`              ]nn$QQQQQQQmmQgmQgQQQWQQQQ5n}`)
	fmt.Println(`               ]{nV$QQQQQQQQQQQQQQQQQQVnn"`)
	fmt.Println(`                 "nnn3HVVHQQQQQWVVVHvnnr'`)
	fmt.Println(`                   "{nnnnvnnnnnnnnnnn"'`)
	fmt.Println(`                      7""nnnnnnn}"""`)
	fmt.Println(`----------------------------------------------------------`)

}

func main(){
	fmt.Println(`----------------------------------------------------------`)
	fmt.Println("            Starting USAGI - Zombifying Daemon")
	printUsagi()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	shutdown := false

	cfg := readConfig()
	cfg = setDefault(cfg)
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
				run := checkIfUp(v.Search)
				if !run {
					//Start process since its not running
					go start(v.Path, v.Param)
				}
			}
			//finished checking, go to sleep
			log.Printf("Check complete, sleeping for %v.\n", cfg.Waittime * time.Millisecond)
			time.Sleep(cfg.Waittime * time.Millisecond)
		}
	}

	fmt.Println("USAGI has shutdown gracefully. Sayonara!")
}
