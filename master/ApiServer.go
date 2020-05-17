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

//post /job/del name=job1
func handlerJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)

	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	name = r.PostForm.Get("name")
	//fmt.Println(name)

	if oldJob, err = G_JobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), oldJob); err == nil {
		w.Write(bytes)
	}
	return
}

func handlerJobWorker(w http.ResponseWriter, r *http.Request) {
	var workerList []string
	var err error
	var bytes []byte

	if workerList, err = G_workerMgr.ListWorkers(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", workerList); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

func handlerJobList(w http.ResponseWriter, r *http.Request) {
	var jobList []*common.Job
	var err error
	var bytes []byte

	if jobList, err = G_JobMgr.ListJob(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		w.Write(bytes)
	}

	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}

}

func handlerJobKill(w http.ResponseWriter, r *http.Request) {
	var err error
	var name string
	var bytes []byte

	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	name = r.PostForm.Get("name")
	if err = G_JobMgr.KillJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

func handlerJobLog(w http.ResponseWriter, r *http.Request) {
	var err error
	var name string
	var skipParam string
	var limitParam string
	var skip int
	var limit int
	var bytes []byte
	var logArr []*common.JobLog

	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	name = r.Form.Get("name")
	skipParam = r.Form.Get("skip")
	limitParam = r.Form.Get("limit")

	if skip, err = strconv.Atoi(skipParam); err != nil {
		skip = 0
	}

	if limit, err = strconv.Atoi(limitParam); err != nil {
		limit = 20
	}

	if logArr, err = G_logMgr.ListLog(name, skip, limit); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", logArr); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

var (
	//单例对象
	GapiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpserver    *http.Server
		staticDir     http.Dir
		staticHandler http.Handler
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handlerJobSave)
	mux.HandleFunc("/job/del", handlerJobDelete)
	mux.HandleFunc("/job/joblist", handlerJobList)
	mux.HandleFunc("/job/jobkill", handlerJobKill)
	mux.HandleFunc("/job/log", handlerJobLog)
	mux.HandleFunc("/worker/list", handlerJobWorker)

	staticDir = http.Dir(G_Config.Webroot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", staticHandler)

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
