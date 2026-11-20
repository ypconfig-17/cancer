package cancer

type Task interface {
	getTopic() string
	getPayload() []byte
}

func (mt *MsgTask) getTopic() string {
	return mt.queueName
}

func (mt *MsgTask) getPayload() []byte {
	return mt.Payload
}

type MsgTask struct {
	Offset    uint64
	Timestamp int64
	Length    int64
	Payload   []byte
	queueName string
	topic     string
}
