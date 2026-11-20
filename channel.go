package cancer

import (
	"cancer/tools"
	"golang.org/x/net/websocket"
	"sync"
	"sync/atomic"
	"time"
)

// ChannelMgr 用以管理channel通道

type Mgr struct {
	sync.Mutex
	channelMap map[string]*Channel
}

var ChannelMgr = &Mgr{
	channelMap: make(map[string]*Channel),
}

func (mgr *Mgr) Create(conn *websocket.Conn) *Channel {
	mgr.Lock()
	defer mgr.Unlock()
	channelMap := mgr.channelMap

	addr := conn.Request().RemoteAddr
	channel := channelMap[addr]
	if channel != nil {
		logger.Warnf("该addr:%s,对应channel已经开启", addr)
		return channel
	}
	channel = newChannel(conn)
	channelMap[addr] = channel
	return channel
}

type Channel struct {
	Name      string
	SessionId string
	Conn      *websocket.Conn
	Address   string
	pipeline  *Pipeline
	status    int32 // 0默认 1开启 2建立连接
}

func newChannel(conn *websocket.Conn) *Channel {
	channel := &Channel{
		Conn:      conn,
		SessionId: tools.TimeUUID().String(),
		Address:   conn.Request().RemoteAddr,
		Name:      "Channel:" + conn.Request().RemoteAddr,
		status:    0,
	}
	channel.pipeline = newPipe(channel)
	return channel
}
func (c *Channel) Start() {
	atomic.CompareAndSwapInt32(&c.status, 0, 1)
	logger.Infof("channel:[%s] 启动成功", c.Name)
	for {
		err := c.receive()
		if err != nil {
			logger.Errorf("channel[%s] receive message error:%s ,will be disconnect", c.Name, err)
			break
		}
	}
}

func (c *Channel) GetPipe() *Pipeline {
	return c.pipeline
}

func (c *Channel) receive() error {

	var reply string
	err := websocket.Message.Receive(c.Conn, &reply)
	if err != nil {
		return err
	}
	c.pipeline.HandlerMsg(reply)
	return nil
}

func (c *Channel) Send(bytes []byte) error {
	_, err := c.Conn.Write(bytes)
	return err
}
func (c *Channel) onClose() {
	if c.closed() {
		return
	}
	defer atomic.CompareAndSwapInt32(&c.status, 1, 0)
	// 休眠300毫秒
	time.Sleep(time.Duration(300) * time.Millisecond)
	// 关闭连接
	err := c.Conn.Close()
	if err != nil {
		logger.Errorf("channel:[%s]异常,error:", c.Name, err)
	}
}

func (c *Channel) onConnect() {
	if !c.opened() {
		logger.Errorf("channel:[%s] 连接建立连接失败，status：%d", c.Name, c.status)
		return
	}
	atomic.CompareAndSwapInt32(&c.status, 1, 2)
	logger.Infof("channel:[%s],建立连接成功,status：%d", c.Name, c.status)

}

func (c *Channel) opened() bool {
	return c.status == 1
}
func (c *Channel) closed() bool {
	return c.status == 0
}
func (c *Channel) connnected() bool {
	return c.status == 2
}
