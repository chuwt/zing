package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	*redis.Client
}

type Config struct {
	Host         string `yaml:"host"`
	Port         int32  `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	DialTimeout  int    `yaml:"dial_timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

func (r *Config) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func NewRedis(cfg Config) *Redis {
	return &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:         cfg.Addr(),
			Password:     cfg.Password, // no password set
			DB:           cfg.DB,       // use default DB
			DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		}),
	}
}
