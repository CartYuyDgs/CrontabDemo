package worker

import (
	"CrontabDemo/common"
	"context"
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
		cmd    *exec.Cmd
		err    error
		outPut []byte
		resutl *common.JobExecuteResult
	)
	go func() {

		resutl = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		resutl.StartTime = time.Now()
		//执行shell
		exec.CommandContext(context.TODO(), "c:\\cygwin64\\bin\\bash.exe", "-c", info.Job.Command)
		outPut, err = cmd.CombinedOutput()

		resutl.EndTime = time.Now()
		//结果返回
		resutl.Output = outPut
		resutl.Err = err

		G_Scheduler.PushJobResult(resutl)

	}()
}

func InitExector() (err error) {
	G_Exector = &Exector{}
	return
}
