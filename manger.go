package cancer

import (
	"cancer/store"
	"sync"
)

const (
	STATE_RUNNABLE     = 0
	STATE_INITIALIZING = 1
	STATE_RUNNING      = 2
	STATE_TERMINATING  = 3
)

// 队列管理员
type Manger struct {
	sync.RWMutex
	status         int8
	fileQueueStore *store.Store
}

func loadAll() {

}

func (mgr *Manger) start() {
	mgr.Lock()
	defer mgr.Unlock()
	mgr.status = 1
}

func (mgr *Manger) stop() {
	mgr.status = 0
}
