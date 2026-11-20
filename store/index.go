package store

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
)

type Idx struct {
	Offset uint64
	LPos   int64
	HPos   int64
}

type Index struct {
	file *os.File
}

func NewIndex(file *os.File) *Index {
	offset := &Index{
		file: file,
	}
	return offset
}

func (index *Index) syncIdx(msgOffset uint64, pos int64, len int64) error {
	idx := &Idx{
		Offset: msgOffset,
		LPos:   pos,
		HPos:   pos + len,
	}
	file := index.file
	buf := idx.encoder()
	_, err := file.Write(buf)
	if err != nil {
		return nil
	}
	return nil
}

func (index *Index) getCurrentMsgPos(offset uint64) int64 {
	// idx 长度24 字节
	if offset > 1 {
		offset = offset - 1
	}
	lPos := 24 * (offset)
	file := index.file

	rBuf := make([]byte, 24)
	_, err := file.ReadAt(rBuf, int64(lPos))
	if err != nil {
		if err == io.EOF {
			return 0
		}
		panic(err)
	}
	buffer := bytes.NewReader(rBuf[:24])
	idx := Idx{}
	err = binary.Read(buffer, binary.BigEndian, &idx)
	if err != nil {
		panic(err)
	}
	return idx.LPos
}

func (idx *Idx) encoder() []byte {
	var buffer bytes.Buffer
	data := Idx{
		Offset: idx.Offset,
		LPos:   idx.LPos,
		HPos:   idx.HPos,
	}

	err := binary.Write(&buffer, binary.BigEndian, &data)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
