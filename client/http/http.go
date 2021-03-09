package http

import (
	"errors"
	"github.com/chuwt/fasthttp-client"
	"github.com/chuwt/zing/json"
	"github.com/chuwt/zing/object"
	"go.uber.org/zap"
)

var log = zap.L().With(zap.Namespace("rest-request"))

// 同步请求
func SyncGetRequest(path string, resp interface{}, mapper object.Params) error {
	client := fasthttp.NewClient()
	if mapper != nil {
		client.AddParams(fasthttp.Mapper(mapper))
	}
	res, err := client.Get(path)
	if err != nil {
		log.Error("同步请求失败", zap.Error(err))
		return err
	}
	log.Debug("同步请求返回", zap.String("body", string(res.Body)))
	if err := json.Json.Unmarshal(res.Body, &resp); err != nil {
		return errors.New(string(res.Body))
	}
	return nil
}

// 异步请求
func AsyncRequest() {

}

type Callback interface {
}
