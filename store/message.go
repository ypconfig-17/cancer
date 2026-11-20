package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

var (
	MsgHeaderLen = int64(24)
)

type Message struct {
	Offset    uint64
	Timestamp int64
	Length    int64
	Payload   []byte
}

func CreteMesage(queue *FileQueue, payload []byte) *Message {
	// 生成下一个offset

	msg := new(Message)
	msg.Timestamp = time.Now().Unix()
	msg.Length = int64(len(payload)) + MsgHeaderLen
	msg.Payload = payload
	return msg
}

type Header struct {
	Offset    uint64
	Timestamp int64
	Length    int64
}

func (msg *Message) encoder() []byte {
	var buffer bytes.Buffer
	// processor header
	header := Header{}
	header.Offset = msg.Offset
	header.Timestamp = msg.Timestamp
	header.Length = msg.Length
	err := binary.Write(&buffer, binary.BigEndian, &header)
	if err != nil {
		log.Panic(err)
	}
	// processor body
	body := msg.Payload
	err = binary.Write(&buffer, binary.BigEndian, &body)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func (msg *Message) decoderHead(buf []byte) {

	header := Header{}

	buffer := bytes.NewReader(buf[:MsgHeaderLen])
	err := binary.Read(buffer, binary.BigEndian, &header)
	if err != nil {
		log.Panic(err)
	}
	msg.Offset = header.Offset
	msg.Timestamp = header.Timestamp
	msg.Length = header.Length

}

func (msg *Message) decoderBody(buf []byte) {
	buffer := bytes.NewReader(buf)

	size := msg.Length - MsgHeaderLen
	body := make([]byte, size)
	//var body [10]byte

	err := binary.Read(buffer, binary.BigEndian, &body)
	if err != nil {
		log.Panic(err)
	}
	msg.Payload = body
}

func (msg *Message) ToString() string {
	toString := fmt.Sprintf("{  Offset: %d,Length:%d , Payload:%s }",
		msg.Offset, msg.Length, string(msg.Payload))
	return toString
}
