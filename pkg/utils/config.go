/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 14:56
@Name: config.go
@Description:
*/

package utils

import (
	"github.com/spf13/viper"
	"strings"
)

var (
	cfgFile = "config"
	cfgDir  = "conf"
)

var v *viper.Viper

func ConfigInit(name, dir string) error {
	v = viper.New()
	if dir != "" {
		v.AddConfigPath(dir)
	} else {
		v.AddConfigPath(cfgDir)
	}
	if name != "" {
		v.SetConfigName(name)
	} else {
		v.SetConfigName(cfgFile)
	}
	v.SetConfigType("yaml")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func GetConfig(c interface{}) error {
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(c); err != nil {
		return err
	}
	return nil
}

func ConfigUtils() *viper.Viper {
	return v
}
