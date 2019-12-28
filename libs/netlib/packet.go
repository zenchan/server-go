package netlib

import (
	"encoding/binary"
	"errors"
)

const (
	defaultConnBufferSize = 4096 // 4KB
)

// packet errors
var (
	ErrIllegalPacket  = errors.New("illegal packet")
	ErrPacketTooLarge = errors.New("packet head too large")
)

// 4 bytes: first two bytes is "ZX"(0x5A58), legal packet flag.
// Last two bytes is packet length.
func checkHead(h []byte) (l int, err error) {
	pf := binary.BigEndian.Uint16(h[:2])
	if pf != 0x5A58 {
		err = ErrIllegalPacket
		return
	}
	n := binary.BigEndian.Uint16(h[2:])
	l = int(n)
	if l > defaultConnBufferSize {
		err = ErrPacketTooLarge
		return
	}
	return
}
