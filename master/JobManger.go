package master

import (
	"CrontabDemo/common"
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
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

func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var jobKey string
	var jobValue []byte
	var putResp *clientv3.PutResponse
	var oldJobObj common.Job

	jobKey = common.JobSaveDir + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	//保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var jobKey string
	var delResp *clientv3.DeleteResponse
	var oldObj common.Job

	jobKey = common.JobSaveDir + name

	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldObj); err != nil {
			return
		}
		oldJob = &oldObj
	}

	return
}

func (jobMgr *JobMgr) ListJob() (jobList []*common.Job, err error) {
	var dirKey string
	var getResp *clientv3.GetResponse
	//var kvPair *mvccpb.KeyValue
	var job *common.Job

	dirKey = common.JobSaveDir

	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}
	jobList = make([]*common.Job, 0)

	for _, kvPair := range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}

		jobList = append(jobList, job)
	}

	return
}

func (jobMgr *JobMgr) KillJob(name string) (err error) {
	var killkey string
	var leaseGrantResp *clientv3.LeaseGrantResponse
	var leaseid clientv3.LeaseID

	killkey = common.JobKill + name

	//让worker监听到一次put操作，自动过期
	if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	leaseid = leaseGrantResp.ID

	if _, err = jobMgr.kv.Put(context.TODO(), killkey, "", clientv3.WithLease(leaseid)); err != nil {
		return
	}

	return
}
