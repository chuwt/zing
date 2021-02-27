package rest

import (
	"github.com/chuwt/zing/gateway"
)

type HuoBi struct {
	Host string
	api  gateway.ApiAuth
}

func NewRest(host string) HuoBi {
	return HuoBi{
		Host: host,
	}
}

func (r *HuoBi) AddAuth(api gateway.ApiAuth) {
	r.api = api
}

func (r *HuoBi) GetUrl(path string) string {
	return r.Host + path
}

type BaseRes struct {
	Status string `json:"status"`
}
