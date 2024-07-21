package app

import (
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
)

type Config struct {
	Server      ServerConfig  `mapstructure:"server"`
	Log         log.Config    `mapstructure:"log"`
	MiddleWare  mid.LogConfig `mapstructure:"middleware"`
	ProjectId   string        `mapstructure:"project_id"`
	Credentials string        `mapstructure:"credentials"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int64  `mapstructure:"port"`
}
