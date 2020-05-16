package worker

import (
	"CrontabDemo/common"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//mongodb存储

type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	G_logSink *LogSink
)

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	//if client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://192.168.31.205:27017")); err != nil{
	//	fmt.Println("++++++++++++",err)
	//	return
	//}
	//fmt.Println(G_Config.MongodbUri)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(G_Config.MongoTimeout)*time.Millisecond)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(G_Config.MongodbUri))
	if err != nil {
		return
	}

	//选择db和collection
	G_logSink = &LogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	go G_logSink.writeLoop()
	return nil
}

func (logSink *LogSink) writeLoop() {
	var (
		log          *common.JobLog
		logBatch     *common.LogBatch
		commitTimer  *time.Timer
		timeoutBatch *common.LogBatch
	)

	for {
		select {
		case log = <-logSink.logChan:
			//写入log
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				commitTimer = time.AfterFunc(time.Duration(G_Config.JobLogCommitTimeout)*time.Millisecond, func(logBatch *common.LogBatch) func() {
					//发出通知
					return func() {
						logSink.autoCommitChan <- logBatch
					}
				}(logBatch))
			}
			logBatch.Logs = append(logBatch.Logs, log)

			if len(logBatch.Logs) >= G_Config.JobLogBatchSize {
				//发送日志
				logSink.saveLogs(logBatch)
				logBatch = nil
				commitTimer.Stop()
			}
		case timeoutBatch = <-logSink.autoCommitChan:
			if timeoutBatch != logBatch {
				continue
			}
			//写入
			logSink.saveLogs(timeoutBatch)
			logBatch = nil
		}
	}
}

func (logSink *LogSink) saveLogs(batch *common.LogBatch) {
	fmt.Println("=======insertmany------------", batch.Logs)
	if _, err := logSink.logCollection.InsertMany(context.TODO(), batch.Logs); err != nil {
		fmt.Println("--------", err)
	}
}

func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:

	}

}
