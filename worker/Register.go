package worker

import (
	"CrontabDemo/common"
	"context"
	"github.com/coreos/etcd/clientv3"
	"net"
	"time"
)

type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIp string
}

var (
	G_register *Register
)

func InitRegister() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   G_Config.EtcdEndPoints,                                     // 集群地址
		DialTimeout: time.Duration(G_Config.EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	if localIp, err = getLocalIp(); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_register = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIp: localIp,
	}

	go G_register.keepOnline()
	return
}

func (register *Register) keepOnline() {
	var (
		regkey         string
		cancleFunc     context.CancelFunc
		cancelCtx      context.Context
		err            error
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
	)

	for {
		regkey = common.Jobworker + register.localIp

		cancleFunc = nil

		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}

		if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			goto RETRY
		}

		cancelCtx, cancleFunc = context.WithCancel(context.TODO())

		if _, err = register.kv.Put(cancelCtx, regkey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}

		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil {
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancleFunc != nil {
			cancleFunc()
		}
	}
}

func getLocalIp() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)

	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}
