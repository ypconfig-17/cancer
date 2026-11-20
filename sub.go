package cancer

import (
	"cancer/stomp"
	"cancer/store"
)

type subscribeMgr struct {
	subscribes map[string][]*Subscribe // key destination value []Subscribe
	store      *store.Store            // key destination value tasklooper
}

func (subMgr *subscribeMgr) SetStore(store *store.Store) {
	subMgr.store = store
}

var subMgr = &subscribeMgr{
	subscribes: make(map[string][]*Subscribe),
}

func OnSubscribe(frame *stomp.Frame, channel *Channel) *stomp.Frame {
	/// validate
	dest := frame.GetHeader(stomp.StompHeaders.Destination)
	subId := frame.GetHeader(stomp.StompHeaders.Id)
	ack := frame.GetHeader(stomp.StompHeaders.Ack)

	var ackFrame *stomp.Frame
	if ack == "client" || ack == "client-individual" {
		ackFrame = stomp.NewAckFrame(subId)
	}

	if dest == "" {
		ackFrame = stomp.NewErrorFrame("Headers missing :'destination'", "Subscribe frame header['destination']: must not be null", frame)
	}
	if subId == "" {
		ackFrame = stomp.NewErrorFrame("Headers missing :'id'", "Subscribe frame header['id']: must not be null", frame)
	}

	subscribe := GetSubscribe(subId, dest)
	if subscribe != nil {
		ackFrame = stomp.NewErrorFrame("Subscribe: ['id'] must uniquely",
			"This sub-id is alread been used with this destination:["+subscribe.id+"].\n sub-id must be uniquely", frame)

	}

	sub := &Subscribe{
		id:          subId,
		destination: dest,
		channel:     channel,
	}
	err := subMgr.registerSubscribe(sub)
	if err != nil {
		ackFrame = stomp.NewErrorFrame("Create subscribe destination error",
			"sever create subscribe:["+subscribe.id+"] error, code 503", frame)

	}
	return ackFrame
}

func (subMgr *subscribeMgr) registerSubscribe(sub *Subscribe) error {

	// 注册消费者
	// 一个destination 多个Subscribe
	subscribes := subMgr.subscribes[sub.destination]
	if subscribes != nil {
		for _, value := range subscribes {
			if value == sub {
				subscribes = append(subscribes, sub)
			}
		}
	} else {

		subscribes = make([]*Subscribe, 10)
		subscribes = append(subscribes, sub)
		subMgr.subscribes[sub.destination] = subscribes
	}
	err := executeSub(sub)
	return err
}

func executeSub(sub *Subscribe) error {
	// 为每一个Subscribe开启一个looper任务
	fileQueue, err := subMgr.store.OpenFileQueue(sub.destination)
	if err != nil {
		return err
	}
	poller := NewTaskPooler(fileQueue)
	looper := NewTaskLooper(10000, sub)
	sub.taskLooper = looper
	driver := NewSubscribeDriver("SubscribeDriver_"+sub.id, poller, looper)
	driver.Start()
	return nil
}

func UnSubscribe(frame *stomp.Frame) *stomp.Frame {
	/// validate
	dest := frame.GetHeader(stomp.StompHeaders.Subscription)
	subId := frame.GetHeader(stomp.StompHeaders.Id)
	ack := frame.GetHeader(stomp.StompHeaders.Ack)

	var ackFrame *stomp.Frame
	if ack == "client" || ack == "client-individual" {
		ackFrame = stomp.NewAckFrame(subId)
	}
	if dest == "" {
		ackFrame = stomp.NewErrorFrame("Headers missing :'destination'", "UnSubscribe frame header['destination']: must not be null", frame)
	}
	if subId == "" {
		ackFrame = stomp.NewErrorFrame("Headers missing :'id'", "UnSubscribe frame header['id']: must not be null", frame)
	}
	delete(subMgr.subscribes, subId)

	return ackFrame
}

func GetSubscribe(subId string, destination string) *Subscribe {
	subscribes := subMgr.subscribes[destination]
	for _, value := range subscribes {
		id := value.id
		if id == subId {
			return value
		}
	}
	return nil
}

type Subscribe struct {
	id          string
	destination string
	channel     *Channel
	taskLooper  *TaskLooper
}

func (sub *Subscribe) GetChannel() *Channel {
	return sub.channel
}
