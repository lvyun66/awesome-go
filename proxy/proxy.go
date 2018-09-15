package proxy

import (
	"encoding/json"
	"github.com/lvyun66/awesome-go/netease/conf"
	"io/ioutil"
	"log"
	"net/http"
)

type Proxy struct {
	ID    int64  `json:"id"`
	IP    string `json:"ip"`
	Type1 string `json:"type1"`
	Speed int    `json:"speed"`
}

func GetProxy() *Proxy {
	response, err := http.Get(conf.DefaultConf.Services.Proxy.Url)
	if err != nil {
		log.Fatalln("Get proxy error, ", err)
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	proxy := &Proxy{}
	json.Unmarshal(data, proxy)
	//if proxy.IP == "" {
	//	log.Fatalln("[PROXY] proxy is empty")
	//}
	return proxy
}
