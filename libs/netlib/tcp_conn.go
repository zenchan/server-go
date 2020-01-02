package netlib

import (
	"io"
	"net"
)

// TCPConn tcp connection
type TCPConn struct {
	net.Conn
	buff []byte

	// readTimeout  time.Duration
	// writeTimeout time.Duration
}

// NewTCPConn news a *TCPConn
func NewTCPConn(c net.Conn) *TCPConn {
	return &TCPConn{
		Conn: c,
		buff: make([]byte, defaultConnBufferSize),
	}
}

// ReadPacket reads packet bytes
func (c *TCPConn) ReadPacket() (buff []byte, err error) {
	var lb [8]byte
	buff = lb[:]
	if _, err = io.ReadFull(c.Conn, buff); err != nil {
		return
	}

	l, err := checkHead(buff)
	if err != nil {
		return
	}

	// auto increse connection's buffer
	if cap(c.buff) < l {
		c.buff = make([]byte, l)
	}

	buff = c.buff[:l]
	if _, err = io.ReadFull(c.Conn, buff); err != nil {
		return
	}
	return
}
