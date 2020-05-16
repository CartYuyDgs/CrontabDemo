package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

//调度计划
type JobSchdulePlan struct {
	Job      *Job
	Expr     *cronexpr.Expression
	NextTime time.Time
}

//执行状态
type JobExecuteInfo struct {
	Job        *Job //任务信息
	PlanTime   time.Time
	RealTime   time.Time
	CancelCtx  context.Context
	CancleFunc context.CancelFunc
}
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo
	Output      []byte
	Err         error
	StartTime   time.Time
	EndTime     time.Time
}

type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`
	Command      string `json:"command" bson:"command"`
	Err          string `json:"err" bson:"err"`
	Output       string `json:"output" bson:"output"`
	PlanTime     int64  `json:"planTime" bson:"planTime"`
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"`
	StartTime    int64  `json:"startTime" bson:"startTime"`
	EndTime      int64  `json:"endTime" bson:"endTime"`
}

//日志批次
type LogBatch struct {
	Logs []interface{}
}

type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
}

//http interface response

type Response struct {
	Erron int         `json:"erron"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(erron int, msg string, data interface{}) (resp []byte, err error) {

	var res Response

	res.Erron = erron
	res.Msg = msg
	res.Data = data

	//序列化
	resp, err = json.Marshal(res)

	return
}

//反序列化
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}

	ret = job
	return
}

func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JobSaveDir)
}

func ExtractJobKillName(killKey string) string {
	return strings.TrimPrefix(killKey, JobKill)
}

//任务变化事件

type JobEvent struct {
	EventType int
	Job       *Job
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildJobSchedulePlan(job *Job) (jobSchedulaPlan *JobSchdulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulaPlan = &JobSchdulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

func BuildJobExecuteInfo(jobScjedulePlan *JobSchdulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobScjedulePlan.Job,
		PlanTime: jobScjedulePlan.NextTime,
		RealTime: time.Now(),
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancleFunc = context.WithCancel(context.TODO())
	return
}
