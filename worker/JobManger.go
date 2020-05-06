package worker

import (
	"CrontabDemo/common"
	"context"
	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_JobMgr *JobMgr
)

//监听任务变化
func (jobMgr *JobMgr) watchJobs(err error) {
	var (
		getReps *clientv3.GetResponse
		kvpair  *mvccpb.KeyValue
		job     *common.Job
	)
	//1.get /cron/jobs目录下的任务，获取当前集群的revision
	if getReps, err = jobMgr.kv.Get(context.TODO(), common.JobSaveDir, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvpair = range getReps.Kvs {
		job, err = common.UnpackJob(kvpair.Value)
		if err == nil {
			//TODO:
		}

	}
	job = job
	return
}

//初始化
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		Kv     clientv3.KV
		lease  clientv3.Lease
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

	G_JobMgr = &JobMgr{
		client: client,
		kv:     Kv,
		lease:  lease,
	}

	return
}
