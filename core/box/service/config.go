package service

import (
	"io/ioutil"
	"os"

	"github.com/ipfs/kubo/repo/fsrepo"
	"github.com/pelletier/go-toml"
)

//Config 为系统全局配置
type ConfigBox struct {
	Sqlite       string `toml:"sqlite"`
	HttpServer   string `toml:"httpServer"`
	IsPrivateNet bool   `toml:"isPrivateNet"`
}

//读取配置文件
func GetConfig() (ConfigBox, error) {
	initConfigBox()
	repoPath, _ := fsrepo.BestKnownPath()
	filePath := repoPath + "/configBox.toml"
	config := ConfigBox{}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Infof("解析config.toml读取错误: %v", err)
		return config, err
	}
	if toml.Unmarshal(content, &config) != nil {
		log.Infof("解析config.toml出错: %v", err)
		return config, err
	}
	return config, err
}

//初始化配置文件
func initConfigBox() {
	repoPath, _ := fsrepo.BestKnownPath()
	filePath := repoPath + "/configBox.toml"
	if IsFileExist(filePath) {
		//log.Info("configBox.toml 已经存在")
	} else {
		configBox := ConfigBox{
			Sqlite:       "box.db",
			HttpServer:   "0.0.0.0:9988",
			IsPrivateNet: true,
		}
		data, err := toml.Marshal(configBox)
		if err != nil {
			log.Infof("解析config.toml出错: %v", err)
			panic("toml error")
		}
		err = ioutil.WriteFile(filePath, data, 0777)
		if err != nil {
			log.Infof("解析config.toml出错: %v", err)
			panic("toml error")
		}
	}

}

func IsFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
