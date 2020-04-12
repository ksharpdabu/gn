package client

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	headerLen = 2
	maxLen    = 1024
)

// Codec 编解码器，用来处理tcp的拆包粘包
type Codec struct {
	Conn    net.Conn
	ReadBuf buffer // 读缓冲
}

// NewCodec 创建一个编解码器
func NewCodec(conn net.Conn) *Codec {
	return &Codec{
		Conn:    conn,
		ReadBuf: newBuffer(make([]byte, 1024)),
	}
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Codec) Read() (int, error) {
	return c.ReadBuf.readFromReader(c.Conn)
}

// Decode 解码数据
// Package 代表一个解码包
// bool 标识是否还有可读数据
func (c *Codec) Decode() ([]byte, bool, error) {
	var err error
	// 读取数据长度
	lenBuf, err := c.ReadBuf.seek(0, headerLen)
	if err != nil {
		return nil, false, nil
	}

	// 读取数据内容
	valueLen := int(binary.BigEndian.Uint16(lenBuf))

	// 数据的字节数组长度大于buffer的长度，返回错误
	if valueLen > maxLen {
		fmt.Println("out of max len")
		return nil, false, nil
	}

	valueBuf, err := c.ReadBuf.read(headerLen, valueLen)
	if err != nil {
		return nil, false, nil
	}
	return valueBuf, true, nil
}
