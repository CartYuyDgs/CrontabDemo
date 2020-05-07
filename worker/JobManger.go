package worker

import (
	"CrontabDemo/common"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

//"go.etcd.io/etcd/mvcc/mvccpb"

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_JobMgr *JobMgr
)

//监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getReps            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)
	//1.get /cron/jobs目录下的任务，获取当前集群的revision
	if getReps, err = jobMgr.kv.Get(context.TODO(), common.JobSaveDir, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvpair = range getReps.Kvs {
		job, err = common.UnpackJob(kvpair.Value)
		if err == nil {
			//TODO:
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			fmt.Println(*jobEvent)
		}

	}
	go func() {
		watchStartRevision = getReps.Header.Revision + 1
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JobSaveDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix()) //监听

		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					//TODO:反序列化 推给scheduler
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
					fmt.Println(*jobEvent)
				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					//TODO: 推一个删除事件给scheduler
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
					fmt.Println(*jobEvent)
				}

			}
		}
	}()
	//job = job
	return
}

//初始化
func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		Kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	//fmt.Printf("*****%s\n",G_Config.EtcdEndPoints)

	config = clientv3.Config{
		Endpoints:   G_Config.EtcdEndPoints,
		DialTimeout: time.Duration(G_Config.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	Kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	G_JobMgr = &JobMgr{
		client:  client,
		kv:      Kv,
		lease:   lease,
		watcher: watcher,
	}

	G_JobMgr.watchJobs()

	return
}
