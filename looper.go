package cancer

import (
	"cancer/stomp"
)

type TaskLooper struct {
	taskQueue chan Task
	cap       int
	sub       *Subscribe
	name      string
}

// 初始化 TaskLooper
func NewTaskLooper(cap int, subscribe *Subscribe) *TaskLooper {
	looper := TaskLooper{
		taskQueue: make(chan Task, cap),
		name:      "TaskLooper-" + subscribe.id,
		sub:       subscribe,
		cap:       cap,
	}
	return &looper
}

func (looper *TaskLooper) overflow() bool {
	if cap(looper.taskQueue) > looper.cap {
		return true
	}
	return false
}

func (looper *TaskLooper) register(task Task) {
	looper.taskQueue <- task
}
func (looper *TaskLooper) batchRegister(taskList []Task) {
	// 当容量小于缓冲容量
	if cap(looper.taskQueue) < looper.cap {
		for _, task := range taskList {
			looper.register(task)
		}
	}
}

func (looper *TaskLooper) loop() {
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("Looper processor message error :%s", err)
				}
			}()
			task := <-looper.taskQueue
			looper.processor(task)
		}()
	}
}

func (looper *TaskLooper) processor(task Task) {
	subscribe := looper.sub
	//pipe := channel.GetPipe()
	payload := task.getPayload()
	topic := task.getTopic()
	frame, _ := stomp.NewMessageFrame(topic, subscribe.id)
	frame.SetBody(payload)
	subscribe.channel.Send([]byte(frame.Serialize()))
	logger.Infof("发送数据:%v", string(payload))
}
