package stomp

import (
	"fmt"
	"strings"
)

// A Frame represents a STOMP frame. A frame consists of a command
// followed by a collection of header entries, and then an optional
// body.
type Frame struct {
	Command string
	Header  *Header
	body    []byte // SetBody
}

func New(command string, headers ...string) *Frame {
	f := &Frame{Command: command, Header: &Header{}}
	for index := 0; index < len(headers); index += 2 {
		f.Header.Add(headers[index], headers[index+1])
	}
	return f
}

// Clone creates a deep copy of the frame and its header. The cloned
// frame shares the body with the original frame.
func (f *Frame) Clone() *Frame {
	fc := &Frame{Command: f.Command}
	if f.Header != nil {
		fc.Header = f.Header.Clone()
	}
	if f.body != nil {
		fc.body = make([]byte, len(f.body))
		copy(fc.body, f.body)
	}
	return fc
}

func (f *Frame) Serialize() string {
	// 数据帧以 COMMAND 命令开始,以 end-of-line (EOL)结束
	frameText := f.Command + "\n"

	for i := 0; i < f.Header.Len(); i++ {
		k, v := f.Header.GetAt(i)
		frameText += k + ":" + v + "\n"
	}
	// header 头信息完毕后,下面跟一个 EOL,区分头部和消息体 Body.
	frameText += "\n"
	frameText += string(f.body)

	// 消息体 Body 后面跟一个 NULL,表示结束
	frameText += "\x00"
	return frameText
}

func (f *Frame) SetBody(body []byte) {
	f.body = body
	f.Header.Add(StompHeaders.ContentLength, fmt.Sprint(len(f.body)))
}

func (f *Frame) Body() []byte {
	return f.body
}

func (f *Frame) GetHeader(key string) string {
	value := f.Header.Get(key)
	return value
}

func Deserialize(stompText string) (*Frame, error) {
	reader := NewReader(strings.NewReader(stompText))
	return reader.Read()
}

func NewAckFrame(id string) *Frame {
	frame := New(CommandACK)
	frame.Header.Add(Id, id)
	//frame.Header.Add(Transaction, id)//事务
	frame.Header.Add(ContentType, Plain+";charset=UTF-8")
	frame.Header.Add(ContentLength, fmt.Sprint(len(frame.body)))
	return frame
}
func NewNAckFrame(id string) *Frame {
	frame := New(CommandACK)
	frame.Header.Add(Id, id)
	//frame.Header.Add(Transaction, id)//事务
	frame.Header.Add(ContentType, Plain+";charset=UTF-8")
	frame.Header.Add(ContentLength, fmt.Sprint(len(frame.body)))
	return frame
}

func NewErrorFrame(tips string, detail string, originFrame *Frame) *Frame {

	frame := New(CommandError)
	if tips != "" {
		frame.Header.Add(Message, tips)
	}
	frame.Header.Add(ContentType, Plain+";charset=UTF-8")
	errorbody := "Server Cannot processed This Frame\n=================\n"

	if originFrame != nil {
		errorbody += originFrame.Serialize()
		errorbody += "\n=================\n"
		if v, isContains := originFrame.Header.Contains(ReceiptId); isContains && v != "" {
			frame.Header.Add(ReceiptId, string(encodeValue(originFrame.Header.Get(ReceiptId))))
		}
	}
	errorbody += detail
	frame.body = []byte(errorbody) //append(frame.Body,errorbody)
	frame.Header.Add(ContentLength, fmt.Sprint(len(frame.body)))

	return frame
}

/*
STOMP 1.2 servers MUST set the following headers:
version : The version of the STOMP protocol the session will be using. See Protocol Negotiation for more details.
STOMP 1.2 servers MAY set the following headers:
heart-beat : The Heart-beating settings.
session : A session identifier that uniquely identifies the session.
server : A field that contains information about the STOMP server. The field MUST contain a server-name field and MAY be followed by optional comment fields delimited by a space character.
The server-name field consists of a name token followed by an optional version number token.
server = name ["/" version] *(comment)
Example:

server:Apache/1.3.9
*/
func NewConnectedFrame(connectFrame *Frame, session string) *Frame {
	connectedFrame := New(FrameC.Connected)
	if connectFrame.Header.Get(StompHeaders.AcceptVersion) != "" {
		maxVersion := -1
		sVerArr := strings.Split(connectFrame.Header.Get(StompHeaders.AcceptVersion), ",")
		for i := len(sVerArr) - 1; i > -1; i-- {
			if sVerArr[i] == "1.2" {
				maxVersion = 2
				break
			} else if sVerArr[i] == "1.1" {
				if maxVersion < 2 {
					maxVersion = 1
				}
			} else if sVerArr[i] == "1.0" {
				if maxVersion < 1 {
					maxVersion = 0
				}
			}
		}
		if maxVersion > -1 {
			connectedFrame.Header.Set(StompHeaders.Version, fmt.Sprint("1.", maxVersion))
		}
		if connectedFrame.Header.Get(StompHeaders.Version) == "" {
			return NewErrorFrame("Version Unsupport", "XStomp Server Supported Protocol Versions Are 1.0,1.1,1.2", connectFrame)
		}
	} else {
		connectedFrame.Header.Set(StompHeaders.Version, "1.0")
	}
	connectedFrame.Header.Set(StompHeaders.Session, session)
	connectedFrame.Header.Set(StompHeaders.Server, "XStompServer/1.0.0")
	connectedFrame.Header.Set(StompHeaders.HeartBeat, "0,0")
	return connectedFrame
}

func NewMessageFrame(destination string, subscription string) (*Frame, string) {
	errorText := ""
	if destination == "" {
		//"The MESSAGE frame MUST include a destination header indicating the destination the message was sent to.
		//If the message has been sent using STOMP,
		//this destination header SHOULD be identical to the one used in the corresponding SEND frame."
		errorText = "destination can not be null/empty." //+
		return nil, errorText
	}
	//if messageId == "" {
	//	errorText = "message-id can not be null/empty." //+ "The MESSAGE frame MUST also contain a message-id header with a unique identifier for that message and a subscription header matching the identifier of the subscription that is receiving the message."
	//	return nil, errorText
	//}

	mHeader := NewHeader()
	mHeader.Add(StompHeaders.Destination, destination)
	//mHeader.Add(StompHeaders.MessageId, messageId)
	if subscription != "" {
		mHeader.Add(StompHeaders.Subscription, subscription)
	}
	mHeader.Add(StompHeaders.ContentType, Plain+";charset=UTF-8")
	messageFrame := &Frame{
		Command: CommandMessage,
		Header:  mHeader,
	}
	return messageFrame, errorText
}
