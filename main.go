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
	"path/filepath"

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
		log.Println(fmt.Sprintf("%s binary located at %s\n\n", b, bin))
		return true
	} else {
		log.Println(fmt.Sprintf("%s binary is not found\n\n", b))
		return false
	}
}

func readConfig(path string) config {
	var c config
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("error opening config file")
	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		log.Println("error unmarshalling config file", err)
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
	args += " &"
	cmd := n + " " + args 
	exec := exec.Command("/bin/sh", "-c", cmd).Run()
	if exec != nil {
		panic(exec)
	}

	log.Printf("Attempting to start %s\n", n)
}

func setDefault(c config) config {
	if c.Waittime == 0 {
		c.Waittime = 60000
	}
	return c
}

func printUsagi(){
	
	log.Println(`----------------------------------------------------------`)
	log.Println(`                g,                        g,`)
	log.Println(`              vQmpg                   _vgQp,`)
	log.Println(`              dQQQmp,                vgWQQQf`)
	log.Println(`             =mQQQQQms             _qWQQQQQ>`)
	log.Println(`              dQQQQQQms           _qWQQQQQQf`)
	log.Println(`              )QQQQQQQms         _qWQQQQQQE'`)
	log.Println(`               mQQQQQQQEggggggggg)QQQQQQQ@f`)
	log.Println(`               ]$QQQQQQQnnnnnnnnnmQQQQQQQf`)
	log.Println(`                ]$QQQQQQEnnnnnnnnQQQQQQQf`)
	log.Println(`                gnVQQQQQmnnnnnnndQQQQQ@vs,`)
	log.Println(`               %nnndQQQQQQQQQQQQQQQQQVvnnn,`)
	log.Println(`              jonnngWQQQQQQQQQQQQQQQQmpvnnn`)
	log.Println(`              oonqmQQ@WQQQQQQQQQQQQ@QQQmnnnL`)
	log.Println(`              onnQQQQQv3H$QQQQQQVVndQQQQEnnc`)
	log.Println(`              nndQQQQEQQnmQQQQQmEmQEQQQQQnn(`)
	log.Println(`              {nmQQQQQggQQQ@VVQQQmgQWQQQQnn'`)
	log.Println(`              ]nn$QQQQQQQmmQgmQgQQQWQQQQ5n}`)
	log.Println(`               ]{nV$QQQQQQQQQQQQQQQQQQVnn"`)
	log.Println(`                 "nnn3HVVHQQQQQWVVVHvnnr'`)
	log.Println(`                   "{nnnnvnnnnnnnnnnn"'`)
	log.Println(`                      7""nnnnnnn}"""`)
	log.Println(`----------------------------------------------------------`)

}

func main(){
	log.Println(`----------------------------------------------------------`)
	log.Println("            Starting USAGI - Zombifying Daemon")
	printUsagi()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Runtime path is: " + dir)
	cfgPath := dir + "/config.yml"
	//shutdown := false
	if len(os.Args) > 1 {
		log.Println("Config file path found in the argument, will proceed with given config file path.")
		cfgPath = os.Args[1]
	}
	log.Println("Reading config file from: " + cfgPath)
	cfg := readConfig(cfgPath)
	cfg = setDefault(cfg)
	for _, v := range cfg.UsagiProspector {
		_, err := exec.LookPath(v.Path)
		if err != nil {
			log.Fatalf("Fatal error: the specified binary for %s cannot be found at path %s", v.Name, v.Path)
	}
	}
	go func(){
		s := <- c
		log.Printf("\nUSAGI is shutting down...\n")
		log.Println(fmt.Sprintf("\n\nReceived signal: %x\n\n", s))
		os.Exit(1)
	}()
	for {
		for _, v := range cfg.UsagiProspector {
			run := checkIfUp(v.Search)
			if !run {
				//Start process since its not running
				start(v.Path, v.Param)
			}
		}
		//finished checking, go to sleep
		log.Println("Check complete, sleeping for 1 minute.")
		time.Sleep(cfg.Waittime * time.Millisecond)
	}

	log.Println("USAGI has shutdown gracefully. Sayonara!")
}
