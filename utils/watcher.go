package utils

import (
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Monitor struct {
	watch *fsnotify.Watcher
}

func SyncFile(WatchPaths []string) {
	M, err := NewMonitor()
	if err != nil {
		log.Println(err)
		return
	}
	defer M.watch.Close()

	for i := 0; i < len(WatchPaths); i++ {
		M.watch.Add(WatchPaths[i])
		log.Println("开始监控：", WatchPaths[i])
	}

	// 退出通道
	done := make(chan bool)
	go M.Do(done)

	// 等待退出信号
	<-done
}

func RmWatcher(watchPath string) {
	M, err := NewMonitor()
	if err != nil {
		log.Println(err)
		return
	}
	defer M.watch.Close()

	M.watch.Remove(watchPath)
	log.Println("移除监控：", watchPath)
}

func NewMonitor() (*Monitor, error) {
	Mon, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Monitor{Mon}, nil
}

func (m *Monitor) Do(done chan bool) {
	defer func() { done <- true }()

	for {
		select {
		case w := <-m.watch.Events:
			// 处理事件
			switch {
			case w.Op&fsnotify.Create == fsnotify.Create:
				fileName := GetFilePath(w.Name)
				if fileIsDir(w.Name) {
					SyncFile([]string{w.Name})
					continue
				}
				time.Sleep(3 * time.Second)
				UploadFile(w.Name, fileName)
				// 处理其他事件（Modify、Delete、Rename等）
			}

		case err := <-m.watch.Errors:
			log.Fatalln(err)
			return
		}
	}
}

// 其他函数保持不变

func fileIsDir(filename string) bool {
	fileHandle, err := os.Stat(filename)
	if err != nil {
		log.Println(err)
		return false
	}
	return fileHandle.IsDir()
}

func GetWatcherPaths(watchPath []string) []string {
	var dirNames []string
	for i := 0; i < len(watchPath); i++ {
		// 读取目录下的子目录
		dirList, err := os.ReadDir(watchPath[i])
		if err != nil {
			log.Println(err)
			continue
		}

		// 监控子目录
		for _, v := range dirList {
			if v.IsDir() {
				dirName := watchPath[i] + v.Name() + string(os.PathSeparator)
				dirNames = append(dirNames, dirName)
				dirNames = append(dirNames, GetWatcherPaths([]string{dirName})...)
			}
		}
	}
	return dirNames
}
