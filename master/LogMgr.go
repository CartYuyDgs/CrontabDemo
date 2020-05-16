package master

import (
	"CrontabDemo/common"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogManger struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogManger
)

func InitLogManger() (err error) {
	var client *mongo.Client
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(G_Config.MongoTimeout)*time.Millisecond)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(G_Config.MongodbUri))
	if err != nil {
		return
	}
	G_logMgr = &LogManger{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

func (logMgr *LogManger) ListLog(name string, skip int, limit int) (logArr []*common.JobLog, err error) {
	var filter *common.JobLogFilter
	var logSort *common.SortLogByStartTime
	var option *options.FindOptions
	var skip_o int64
	var limit_o int64
	var cursor *mongo.Cursor
	var jobLog *common.JobLog

	logArr = make([]*common.JobLog, 0)

	skip_o = int64(skip)
	limit_o = int64(limit)

	filter = &common.JobLogFilter{
		JobName: name,
	}

	logSort = &common.SortLogByStartTime{SortOrder: -1}
	option = &options.FindOptions{
		Limit: &limit_o,
		Skip:  &skip_o,
		Sort:  logSort,
	}

	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter, option); err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		if err = cursor.Decode(jobLog); err != nil {
			continue
		}
		logArr = append(logArr, jobLog)
	}
	return
}
