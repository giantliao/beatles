package streamserver

import (
	"github.com/giantliao/beatles/config"
	"github.com/giantliao/beatles/port"
	"log"
	"net"
	"strconv"
)

type StreamServer struct {
	addr string
	quit chan int
}

func NewStreamServer() *StreamServer {
	tcpport := port.TcpPort()

	addr := ":" + strconv.Itoa(tcpport)

	cfg := config.GetCBtl()
	cfg.SetStreamPort(tcpport)

	return &StreamServer{addr: addr}
}

func (ss *StreamServer) StartServer() error {
	lis, err := net.Listen("tcp", ss.addr)
	if err != nil {
		panic("failed to listen on %s" + ss.addr + " : " + err.Error())
	}

	log.Println("Stream Server start at ", ss.addr)

	return nil
}

func (ss *StreamServer) StopServer() {

}
