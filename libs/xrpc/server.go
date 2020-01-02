package xrpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/zenchan/server-go/libs/xrpc/pb"
	"google.golang.org/grpc"
)

// RequestCallback request/reply callback function
type RequestCallback func([]byte) []byte

// RouteCallback route callback function
type RouteCallback func([]byte) []byte

type rpcServer struct {
	requestCb RequestCallback
	routeCb   RouteCallback
	errCb     func(string)
}

func (s *rpcServer) setDefaultOption() {
	s.requestCb = func(in []byte) []byte { return []byte{} }
	s.routeCb = func(in []byte) []byte { return []byte{} }
	s.errCb = func(string) {}
}

func (s *rpcServer) Request(ctx context.Context, in *pb.LANPacket) (out *pb.LANPacket, err error) {
	buff := s.requestCb(in.Body)
	out = &pb.LANPacket{
		Body: buff,
	}
	return
}

func (s *rpcServer) Route(ins pb.RPC_RouteServer) error {
	for {
		pkt, err := ins.Recv()
		if err == io.EOF {
		}
		if err != nil {
			s.errCb(err.Error())
			return err
		}
		s.routeCb(pkt.Body)
	}
}

var rpcSer = rpcServer{}

// Serve rpc serve
func Serve(port int, opts ...Option) (err error) {
	rpcSer.setDefaultOption()
	for _, opt := range opts {
		opt(&rpcSer)
	}

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	gs := grpc.NewServer()
	pb.RegisterRPCServer(gs, &rpcSer)
	if err = gs.Serve(lis); err != nil {
		log.Printf("grpc serve failed: %s\n", err.Error())
		return
	}

	return
}
