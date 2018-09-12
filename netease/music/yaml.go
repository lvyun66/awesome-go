package music

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
	yamlFile, err := ioutil.ReadFile("./music/app.yml")
	if err != nil {
		log.Fatalln("Loading yaml config error:", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return conf
}
