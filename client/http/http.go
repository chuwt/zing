package http

import (
	"encoding/json"
	"github.com/chuwt/fasthttp-client"
	"go.uber.org/zap"
)

var log = zap.L().With(zap.Namespace("rest-request"))

// 同步请求
func SyncGetRequest(path string, resp interface{}, mapper fasthttp.Mapper) error {
	client := fasthttp.NewClient()
	if mapper != nil {
		client.AddParams(mapper)
	}
	res, err := client.Get(path)
	if err != nil {
		log.Error("同步请求失败", zap.Error(err))
		return err
	}
	log.Debug("同步请求返回",
		zap.String("body", string(res.Body)))
	if err := json.Unmarshal(res.Body, resp); err != nil {
		return err
	}
	return nil
}

// 异步请求
func AsyncRequest() {

}

type Callback interface {
}
