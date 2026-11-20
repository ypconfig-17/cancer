package store

import (
	"cancer/log"
	"errors"
	"fmt"
	"os"
)

var (
	notDirError = errors.New(" opened path is not a directory")
)

var logger = log.RegisterScope("/store")

//var MessageStore, _ = OpenStore("/db", false)

type Store struct {
	dir      string
	readOnly bool
	queues   map[string]*FileQueue //指向具体的数据
}

func (st *Store) OpenFileQueue(queueName string) (*FileQueue, error) {

	fileQueue := st.queues[queueName]
	if fileQueue != nil {
		return fileQueue, nil
	}

	fileQueue, err := openFileQueue(st.dir, queueName)
	if err != nil {
		return nil, fmt.Errorf("create the queue :%s fail,error:%s", queueName, err)
	}
	st.queues[queueName] = fileQueue

	return fileQueue, nil
}

func (st *Store) GetQueue(queueName string) *FileQueue {
	queue := st.queues[queueName]
	return queue
}

func OpenStore(dirPath string, readOnly bool) (*Store, error) {
	fi, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				return nil, err
			}
		}
	}
	if fi == nil {
		fi, err = os.Stat(dirPath)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("opened path:%s is not a directory", dirPath)
	}

	//查看队列
	queuemap := make(map[string]*FileQueue)
	//rd, err := ioutil.ReadDir(dirPath)
	//for _, fi := range rd {
	//	if !fi.IsDir() {
	//		continue
	//	}
	//	if strings.Contains(fi.Name(), "store") || strings.Contains(fi.Name(), "index") ||
	//		strings.Contains(fi.Name(), "offset") || strings.Contains(fi.Name(), "idx") {
	//		continue
	//	}
	//	queueName := fi.Name()
	//	fileQueue, err := openFileQueue(dirPath, queueName)
	//	if err != nil {
	//		return nil, err
	//	}
	//	queuemap[queueName] = fileQueue
	//}
	store := &Store{}
	store.dir = dirPath
	store.queues = queuemap
	store.readOnly = readOnly

	return store, nil

}
