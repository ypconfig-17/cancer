package store

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type Offset struct {
	file *os.File
}

func newOffset(file *os.File) *Offset {
	offset := &Offset{
		file: file,
	}
	return offset
}

func (off *Offset) setOffset(offset uint64) {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, &offset)
	if err != nil {
		panic(err)
	}

	wbuf := buffer.Bytes()
	_, err = off.file.WriteAt(wbuf, 0)
	if err != nil {
		panic(err)
	}
}

func (off *Offset) getCurrentOffset() uint64 {
	var offset uint64
	rBuf := make([]byte, 8)
	_, err := off.file.ReadAt(rBuf, 0)
	if err != nil {
		if err == io.EOF {
			off.setOffset(0)
			return 0
		}
		panic(err)
	}
	buffer := bytes.NewReader(rBuf)
	err = binary.Read(buffer, binary.BigEndian, &offset)
	if err != nil {
		panic(err)
	}

	return offset
}
func (off *Offset) setNext() uint64 {
	offset := off.getCurrentOffset()
	next := offset + 1
	off.setOffset(next)
	return next
}
