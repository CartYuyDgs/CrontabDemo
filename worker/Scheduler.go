package worker

import (
	"CrontabDemo/common"
	"fmt"
	"time"
)

type Scheduler struct {
	jobEventChan      chan *common.JobEvent             //事件队列
	jobPlanTable      map[string]*common.JobSchdulePlan //调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo
}

var G_Scheduler *Scheduler

func InitScheduler() (err error) {
	G_Scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchdulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
	}
	go G_Scheduler.schedulerLoop()
	return nil
}
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchdulePlan) {
	var (
		jobexecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)
	//	调度，执行

	//执行的任务可能执行很久，有可能执行一次 防止并发
	if jobexecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		fmt.Println("执行中......")
		return
	}
	//构建执行状态
	jobexecuteInfo = common.BuildJobExecuteInfo(jobPlan)
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobexecuteInfo

	//执行任务
	//TODO:
	fmt.Println("执行任务：", jobexecuteInfo.Job.Name, jobexecuteInfo.RealTime, jobexecuteInfo.PlanTime)

	return
}

//重新计算任务调度状态
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchdulePlan
		nearTime *time.Time
	)

	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	now := time.Now()
	//遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//尝试执行任务
			scheduler.TryStartJob(jobPlan)
			fmt.Println("执行任务：" + jobPlan.Job.Name)
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}

		//统计最近一个要执行的任务到期时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	//过期的任务立即执行

	//统计最近要过期的任务时间
	scheduleAfter = (*nearTime).Sub(now)
	return
}

//处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchdulePlan
		jobExited       bool
		err             error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan

	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExited = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExited {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}

func (scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)
	scheduleAfter = scheduler.TrySchedule()

	scheduleTimer = time.NewTimer(scheduleAfter)

	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C:

		}
		scheduleAfter = scheduler.TrySchedule()
		scheduleTimer.Reset(scheduleAfter)
	}
}

//推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}
