package utils

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/text/gstr"
)

var configPath string

/**
自动获取环境变量
如果环境变量不存在，获取本地config-local.yaml
*/
func Setup() {
	env := gstr.ToUpper(genv.Get("ENV"))
	configFile := getEnvConfig()
	if env != "" {
		g.Log().SetPrefix("[" + env + "]")
		if value, ok := configFile[env]; ok {
			configPath = value
			g.Log().Info("读取配置文件: %v", configPath)
			g.Cfg().SetFileName(configPath)
		} else {
			g.Log().Error("error: 环境变量未设置或设置错误")
		}
	} else {
		g.Cfg().SetFileName("config-local.yaml")
		g.Log().SetPrefix("[local]")
		g.Log().Info("环境变量未配置,读取配置文件: config-local.yaml")
	}
}

func getEnvConfig() map[string]string {
	return map[string]string{
		"DEV": "config-dev.yaml", //开发环境
		"FAT": "config-fat.yaml", //测试环境
		"UAT": "config-uat.yaml", //预发布环境（灰度环境）
		"PRO": "config-pro.yaml", //生产环境
	}
}
