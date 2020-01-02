package main

import (
	"net"
	"time"

	"github.com/zenchan/server-go/libs/netlib"
	"github.com/zenchan/server-go/libs/xlog"
	"github.com/zenchan/server-go/proto/pb"
)

func tcpLoop(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			xlog.Errorf("tcp accept failed: %s", err.Error())
			continue
		}

		go handleTCPConn(conn)
	}
}

func handleTCPConn(conn net.Conn) {
	tc := netlib.NewTCPConn(conn)
	defer func() {
		tc.Close()
	}()

	tc.SetReadDeadline(time.Now().Add(time.Second * 3))
	buff, err := tc.ReadPacket()
	if err != nil {
		xlog.Infof("conn read first packet failed: %s", err.Error())
		return
	}
	pkt := pb.Packet{}
	if err = pkt.Unmarshal(buff); err != nil {
		xlog.Infof("packet decode failed: %s", err.Error())
		return
	}

	// TODO send to lobby

	for {
		tc.SetReadDeadline(time.Now().Add(time.Minute))
		buff, err := tc.ReadPacket()
		if err != nil {
			xlog.Infof("conn read packet failed: %s", err.Error())
			return
		}

		pkt := pb.Packet{}
		if err = pkt.Unmarshal(buff); err != nil {
			xlog.Infof("packet decode failed: %s", err.Error())
			return
		}
	}
}
