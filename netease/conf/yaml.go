package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Yaml struct {
	Services struct {
		Mysql struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"mysql"`
		Proxy struct {
			Url string `yaml:"url"`
		} `yaml:"proxy"`
	} `yaml:"services"`
}

func LoadConf() *Yaml {
	conf := &Yaml{}
	filePath := os.Getenv("GOPATH") + "/src/github.com/lvyun66/awesome-go/netease/conf/app.yml"
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Loading yaml config error:", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return conf
}

var DefaultConf = LoadConf()
