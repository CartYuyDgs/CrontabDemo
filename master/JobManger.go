package master

import (
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
