package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

var (
	Config = New()
)

/**
 * 返回单例实例
 * @method New
 */
func New() *toml.Tree {
	config, err := toml.LoadFile("./config/config.toml")
	if err != nil {
		fmt.Println("TomlError ", err.Error())
	}
	return config
}
