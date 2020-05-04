package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

func handlerJobSave(w http.ResponseWriter, r *http.Request) {

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
