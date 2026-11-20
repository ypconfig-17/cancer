package cancer

import (
	"cancer/stomp"
	"container/list"
)

type Handler func(cxt *HandlerContext, msg Msg)
type Msg interface{}

type HandlerContext struct {
	channel  *Channel
	pipeline *Pipeline
	current  *list.Element
}

func (cxt *HandlerContext) OnConnected(messageFrame *stomp.Frame) {
	err := stomp.ValidateServerCommand(messageFrame.Command)
	if err != nil {
		logger.Errorf("accept message error:%s:", err)
		return
	}
	stringMsg := messageFrame.Serialize()
	cxt.channel.onConnect()
	cxt.channel.Send([]byte(stringMsg))
}

func (cxt *HandlerContext) OnAck(messageFrame *stomp.Frame) {
	if messageFrame == nil {
		return
	}
	stringMsg := messageFrame.Serialize()
	cxt.channel.Send([]byte(stringMsg))
}

func (cxt *HandlerContext) fireHandler(msg Msg) {
	next := cxt.current.Next()
	if next == nil {
		return
	}
	cxt.current = next

	handler := cxt.current.Value.(Handler)
	handler(cxt, msg)
}

func StompEncoderHandler(cxt *HandlerContext, msg Msg) {
	frame, err := stomp.Deserialize(msg.(string))
	if err != nil {
		logger.Errorf("channel read a message of frame error :{%s}", err)
	}
	cxt.fireHandler(frame)
}

func FrameDecoderHandler(cxt *HandlerContext, msg Msg) {

	frame := msg.(*stomp.Frame)
	channel := cxt.channel

	sessionId := channel.SessionId
	switch frame.Command {
	case stomp.FrameC.Connect:
		/*
		 * 这里需要做一些协议支持处理 http://stomp.github.io/stomp-specification-1.1.html#Protocol_Negotiation
		 * 可接受的协议 Done.
		 * host验证
		 * login
		 * passcode
		 */
		connectedFrame := stomp.NewConnectedFrame(frame, sessionId)
		cxt.OnConnected(connectedFrame)
		break
	case stomp.FrameC.Subscribe:
		returnFrame := OnSubscribe(frame, channel)
		cxt.OnAck(returnFrame)
		break
	case stomp.FrameC.Send:
		// 解析数据
		cxt.fireHandler(msg)
		break
	case stomp.FrameC.Unsubscribe:
		returnFrame := UnSubscribe(frame)
		cxt.OnAck(returnFrame)
		break
	default:
		logger.Warn("server can not support data frame on this command:%s", frame.Command)
		break

	}
	return
}

func MessageEnqueueHandler(cxt *HandlerContext, msg Msg) {
	frame := msg.(*stomp.Frame)
	desc := frame.GetHeader(stomp.StompHeaders.Destination)
	body := frame.Body()
	message := NewPubMessage(desc, body)
	OnMessage(message)
}

func StompDecoderHandler(cxt *HandlerContext, msg Msg) {
	frame := msg.(stomp.Frame)
	msg, err := stomp.Deserialize(frame.Serialize())
	if err != nil {
		logger.Errorf("channel read a message of frame error :{%s}", err)
	}
	cxt.fireHandler(msg)

}
