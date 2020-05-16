package worker

import (
	"CrontabDemo/common"
	"os/exec"
	"time"
)

type Exector struct {
}

var (
	G_Exector *Exector
)

func (exector *Exector) ExectorJob(info *common.JobExecuteInfo) {
	var (
		cmd     *exec.Cmd
		err     error
		outPut  []byte
		resutl  *common.JobExecuteResult
		jobLock *JobLock
	)
	go func() {

		resutl = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		//获取分布式锁
		jobLock = G_JobMgr.CreateJobLock(info.Job.Name)

		resutl.StartTime = time.Now()
		//执行shell
		time.Sleep(time.Duration(1000) * time.Millisecond)
		err = jobLock.TryLock()
		defer jobLock.UnLock()

		if err != nil {
			resutl.Err = err
			resutl.EndTime = time.Now()
		} else {
			resutl.StartTime = time.Now()
			cmd = exec.CommandContext(info.CancelCtx, "c:\\cygwin64\\bin\\bash.exe", "-c", info.Job.Command)
			outPut, err = cmd.CombinedOutput()

			resutl.EndTime = time.Now()
			//结果返回
			resutl.Output = outPut
			resutl.Err = err
		}

		G_Scheduler.PushJobResult(resutl)

	}()
}

func InitExector() (err error) {
	G_Exector = &Exector{}
	return
}
