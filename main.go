package main

import (
	"strings"
	"tocn_api/qiniu-auto-sync/utils"

	"github.com/qiniu/x/log"
)

func main() {
	WatchePath, err := utils.GetOption("WatcherPath", "watcher")

	if err != nil {
		log.Fatalf("请先设置监控目录 %v", err)
		return
	}

	// 配置文件监控的子目录
	dirList := utils.GetWatcherPaths(strings.Split(WatchePath, ";"))
	// 配置文件监控目录
	dirList = append(dirList, strings.Split(WatchePath, ";")...)
	utils.SyncFile(dirList)
}
