package netframe

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
	var lb [4]byte
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

// // ReadHead reads packet head bytes
// func (c *TCPConn) ReadHead() (buff []byte, err error) {
// 	var lb [4]byte
// 	buff = lb[:]
// 	if _, err = io.ReadFull(c.Conn, buff); err != nil {
// 		return
// 	}

// 	hl, err := checkHead(buff)
// 	if err != nil {
// 		return
// 	}

// 	buff = c.buff[:hl]
// 	if _, err = io.ReadFull(c.Conn, buff); err != nil {
// 		return
// 	}
// 	return
// }

// // ReadBody reads packet body bytes
// func (c *TCPConn) ReadBody(ln int) (buff []byte, err error) {
// 	// auto increse connection's buffer
// 	if cap(c.buff) < ln {
// 		c.buff = make([]byte, ln)
// 	}

// 	buff = c.buff[:ln]
// 	if _, err = io.ReadFull(c.Conn, buff); err != nil {
// 		return
// 	}
// 	return
// }

// // SetTimeout set read and write timeout
// func (c *TCPConn) SetTimeout(t time.Duration) {
// 	c.readTimeout = t
// 	c.writeTimeout = t
// }

// // SetReadTimeout set read timeout
// func (c *TCPConn) SetReadTimeout(t time.Duration) {
// 	c.readTimeout = t
// }

// // SetWriteTimeout set write timeout
// func (c *TCPConn) SetWriteTimeout(t time.Duration) {
// 	c.writeTimeout = t
// }
