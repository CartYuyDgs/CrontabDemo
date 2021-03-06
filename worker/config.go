package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdEndPoints       []string `json:"etcdEndPoints"`
	EtcdDialTimeout     int      `json:"etcdDialTimeout"`
	MongoTimeout        int      `json:"mongoTimeout"`
	MongodbUri          string   `json:"mongodbUri"`
	JobLogBatchSize     int      `json:"jobLogBatchSize"`
	JobLogCommitTimeout int      `json:"jobLogCommitTimeout"`
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
