package cancer

import (
	"cancer/store"
	"time"
)

type TaskPoller interface {
	poll() Task
}

type FileQueuePoller struct {
	queue store.Queue
}

func NewTaskPooler(queue store.Queue) TaskPoller {
	poller := &FileQueuePoller{
		queue: queue,
	}
	return poller
}

func (poller *FileQueuePoller) poll() Task {
	message := poller.queue.PollTime(time.Second * 1)
	if message == nil {
		return nil
	}
	task := poller.convertor(message)
	return task
}

func (poller *FileQueuePoller) convertor(message *store.Message) Task {
	queue := poller.queue
	var task Task
	task = &MsgTask{
		Offset:    message.Offset,
		Timestamp: message.Timestamp,
		Length:    message.Length,
		Payload:   message.Payload,
		queueName: queue.Name(),
		topic:     queue.Topic(),
	}
	return task

}
