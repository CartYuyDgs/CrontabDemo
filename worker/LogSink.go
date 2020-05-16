package worker

import (
	"CrontabDemo/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//mongodb存储

type LogSink struct {
	client        *mongo.Client
	logCollection *mongo.Collection
	logChan       chan *common.JobLog
}

var (
	G_logSink *LogSink
)

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	client, err = mongo.NewClient(options.Client().ApplyURI(G_Config.MongodbUri))
	if err != nil {
		return
	}

	//选择db和collection
	G_logSink = &LogSink{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
		logChan:       make(chan *common.JobLog, 1000),
	}

	go G_logSink.writeLoop()

}

func (logSink *LogSink) writeLoop() {
	var (
		log *common.JobLog
	)

	for {
		select {
		case log = <-logSink.logChan:
			//写入log
		}
	}
}
