package zing

import (
	"github.com/chuwt/zing/client/redis"
	"github.com/chuwt/zing/config"
	dataCenter "github.com/chuwt/zing/data_center"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"sync"
)

var (
	Log = zap.L().With(zap.Namespace("zing"))
)

type Zing struct {
	single singleflight.Group
	mu     sync.Mutex

	//pyEngine python.PyEngine

	userCtx UserCtx // 用户上下文
}

func NewZing(redisCfg redis.Config, strategyCfg config.Strategy) Zing {
	return Zing{
		single: singleflight.Group{},
		mu:     sync.Mutex{},
		//pyEngine: python.NewPyEngine(strategyCfg.Path, strategyCfg.PythonPath),

		userCtx: NewUserCtx(dataCenter.NewRedisTower(redisCfg)),
	}
}
