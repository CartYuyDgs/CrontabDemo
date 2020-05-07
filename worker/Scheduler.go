package worker

import "CrontabDemo/common"

type Scheduler struct {
	jobEventChan chan *common.JobEvent             //事件队列
	jobPlanTable map[string]*common.JobSchdulePlan //调度计划表
}

var G_Scheduler *Scheduler

func initScheduler() {
	G_Scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
	}
	go G_Scheduler.schedulerLoop()
	return
}

//处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
	case common.JOB_EVENT_DELETE:

	}
}

func (scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent *common.JobEvent
	)
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			scheduler.handleJobEvent(jobEvent)
		}
	}
}

//推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}
