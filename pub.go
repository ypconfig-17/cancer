package cancer

import (
	"cancer/store"
)

type publishMgr struct {
	store *store.Store // key destination value tasklooper
}

func (pubMgr *publishMgr) SetStore(store *store.Store) {
	pubMgr.store = store
}

// store.Message wapper
type PubMessage struct {
	pubDesc string
	body    []byte // SetBody
}

func NewPubMessage(desc string, body []byte) *PubMessage {
	return &PubMessage{
		pubDesc: desc,
		body:    body,
	}
}

var pubMgr = &publishMgr{}

func OnMessage(pubMsg *PubMessage) {
	err := pubMgr.publish(pubMsg)
	if err != nil {
		logger.Errorf("pubMsg add message store error:%s ", err)
		panic(err)
	}
}

func (pubMgr *publishMgr) publish(pubMsg *PubMessage) error {
	queueStroe := pubMgr.store
	desc := pubMsg.pubDesc
	queue, err := queueStroe.OpenFileQueue(desc)
	if err != nil {
		return err
	}

	// 队列自增offset
	mesage := store.CreteMesage(queue, pubMsg.body)

	queue.Add(mesage)
	return nil
}
