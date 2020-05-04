package master

import (
	"CrontabDemo/common"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

func handlerJobSave(w http.ResponseWriter, r *http.Request) {
	//解析POST表单
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	postJob = r.PostForm.Get("job")

	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	if oldJob, err = G_JobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//	应答返回  ({"error":0，”msg“:"", "data":{......}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
	fmt.Println(err)
}

var (
	//单例对象
	GapiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpserver *http.Server
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handlerJobSave)

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_Config.ApiPort)); err != nil {
		return
	}

	httpserver = &http.Server{
		ReadTimeout:  time.Duration(G_Config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_Config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	GapiServer = &ApiServer{
		httpServer: httpserver,
	}

	go httpserver.Serve(listener)

	return
}
