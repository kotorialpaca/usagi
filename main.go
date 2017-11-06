package main

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
	"gopkg.in/yaml.v2"
)

type config struct {
	UsagiProspector []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
		Param string `yaml:"param"`
	} `yaml:"UsagiProspector"`
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
	err = yaml.Unmarshal(f, c)
	if err != nil {
		fmt.Println("error unmarshalling config file")
	}
	return c
}

func main(){
	fmt.Println("Starting USAGI - Zombifying Daemon")
}
