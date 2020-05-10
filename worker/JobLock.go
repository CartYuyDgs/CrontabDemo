package worker

import (
	"CrontabDemo/common"
	"context"
	"github.com/coreos/etcd/clientv3"
)

type JobLock struct {
	kv        clientv3.KV
	lease     clientv3.Lease
	jobName   string
	canceFunc context.CancelFunc
	leaseId   clientv3.LeaseID
	isLocked  bool
}

func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		jobName: jobName,
		kv:      kv,
		lease:   lease,
	}
	return
}

func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrandResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		canceFunc      context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)
	//创建租约
	if leaseGrandResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	//用于取消续租
	cancelCtx, canceFunc = context.WithCancel(context.TODO())

	leaseId = leaseGrandResp.ID
	//自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}
	//处理续租
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)

		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()
	//创建事务
	txn = jobLock.kv.Txn(context.TODO())

	lockKey = common.JobLockDir + jobLock.jobName
	//事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	//成功返回，失败释放租约
	if !txnResp.Succeeded {
		err = common.ERR_LOCK_ALREADER_REQUIRED
		goto FAIL
	}
	jobLock.isLocked = true
	jobLock.leaseId = leaseId
	jobLock.canceFunc = canceFunc
	return

	return

FAIL:
	canceFunc()
	jobLock.lease.Revoke(context.TODO(), leaseId)
	return
}

func (jobLock *JobLock) UnLock() {
	if jobLock.isLocked {
		jobLock.canceFunc()
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseId)
	}
}
