package common

import "encoding/json"

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

//http interface response

type Response struct {
	Erron int         `json:"erron"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(erron int, msg string, data interface{}) (resp []byte, err error) {

	var res Response

	res.Erron = erron
	res.Msg = msg
	res.Data = data

	//序列化
	resp, err = json.Marshal(res)

	return
}

//反序列化
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}

	ret = job
	return
}
