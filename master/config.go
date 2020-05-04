package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeout  int      `json:"apiReadTimeout"`
	ApiWriteTimeout int      `json:"apiWriteTimeout"`
	EtcdEndPoints   []string `json:"etcdEndPoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}

var (
	G_Config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	//反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_Config = &conf

	//fmt.Println(G_Config)
	return
}
