package store

import (
	"cancer/tools"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Queue interface {
	Add(message *Message) bool
	Poll() *Message
	PollTime(time time.Duration) *Message
	//Seek(offset int)
	//Size() uint64
	//Empty() bool
	Topic() string
	Name() string
}

func (queue *FileQueue) Topic() string {
	return queue.topic
}
func (queue *FileQueue) Name() string {
	return queue.name
}

type FileQueue struct {
	name     string
	topic    string
	dir      string
	path     string
	offset   *Offset
	writePos int64 //当前队列读取的偏移量
	readPos  int64 //当前队列读取的偏移量
	file     *os.File
	index    *Index
	cond     *tools.Condition
}

func (queue *FileQueue) Add(message *Message) bool {
	// processor offset

	nextOffset := queue.offset.setNext()
	message.Offset = nextOffset
	writer := queue.file
	stat, err := writer.Stat()
	if err != nil {
		panic(err)
	}
	size := stat.Size()

	buf := message.encoder()
	msgLen, err := writer.Write(buf)
	if err != nil {
		return false
	}
	//len the message bytes len
	err = queue.index.syncIdx(message.Offset, size, int64(msgLen))
	if err != nil {
		return false
	}
	return true
}

func (queue *FileQueue) PollTime(time time.Duration) *Message {
	message := queue.Poll()
	// 读到文件末尾
	if message == nil {
		queue.cond.AwaitTime(time)
	}
	return message
}

func (queue *FileQueue) Poll() *Message {
	reader := queue.file
	readOff := queue.readPos
	//logger.Infof("当前读取文件offset:%d ", readOff)
	hBuf := make([]byte, MsgHeaderLen)
	_, err := reader.ReadAt(hBuf, readOff)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return nil
	}
	// body  offset = msglen-headerlen

	message := Message{}
	message.decoderHead(hBuf)

	bodyLen := message.Length - MsgHeaderLen
	readOff = readOff + MsgHeaderLen

	bBuf := make([]byte, bodyLen)
	bLen, err := reader.ReadAt(bBuf, readOff)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return nil
	}

	message.decoderBody(bBuf)
	readOff = readOff + int64(bLen)
	queue.readPos = readOff
	return &message

}

func openFileQueue(dirPath string, name string) (*FileQueue, error) {
	// init queue
	filePath := filepath.Join(dirPath, name)

	if strings.Contains(filePath, "/") {
		index := strings.LastIndex(filePath, "/")
		dir := filePath[0:index]
		exist := tools.IsExist(dir)
		if !exist {
			err := tools.MkdirAll(filePath, 0755)
			if err != nil {
				return nil, err
			}
		}
	}

	fileName := filePath + ".log"
	// file
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	// idx
	idxFileName := filePath + "_idx.log"
	indexFile, err := os.OpenFile(idxFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	index := NewIndex(indexFile)
	// offset
	offFileName := filePath + "_offset.log"
	offsetFile, err := os.OpenFile(offFileName, os.O_RDWR|os.O_CREATE, 0664)

	offset := newOffset(offsetFile)
	currentOffset := offset.getCurrentOffset()
	if err != nil {
		return nil, err
	}

	writePos := index.getCurrentMsgPos(currentOffset)

	queue := &FileQueue{
		dir:      dirPath,
		path:     filePath,
		topic:    name,
		name:     name,
		file:     file,
		index:    index,
		offset:   offset,
		writePos: writePos,
		readPos:  0,
		cond:     tools.NewCondition(),
	}

	return queue, nil
}
