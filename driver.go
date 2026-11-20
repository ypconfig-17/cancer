package cancer

import (
	"sync"
	"sync/atomic"
)

type SubscribeDriver struct {
	sync.Mutex
	name   string
	status int32
	poller TaskPoller
	looper *TaskLooper
}

func NewSubscribeDriver(name string, poller TaskPoller, looper *TaskLooper) *SubscribeDriver {
	return &SubscribeDriver{
		poller: poller,
		looper: looper,
		status: 0,
		name:   name,
	}
}

func (sd *SubscribeDriver) Start() {
	swapped := atomic.CompareAndSwapInt32(&sd.status, 0, 1)
	if !swapped {
		logger.Warn("SubscribeDriver has been started")
		return
	}
	logger.Infof("SubscribeDriver will be start")
	go sd.run()
}
func (td *SubscribeDriver) run() {
	looper := td.looper
	poller := td.poller
	go looper.loop()

	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("Looper loop message error :%s", err)
				}
			}()
			// 当容量小于缓冲容量
			if !looper.overflow() {
				task := poller.poll()
				if task != nil {
					looper.register(task)
				}
			}
		}()
	}

}
