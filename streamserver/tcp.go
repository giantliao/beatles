package streamserver

import (
	"github.com/giantliao/beatles-protocol/stream"
	"github.com/giantliao/beatles/config"
	"github.com/giantliao/beatles/port"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type StreamServer struct {
	addr    string
	quit    chan struct{}
	lis     net.Listener
	session map[string]net.Conn
	wg      sync.WaitGroup
}

type CloseConn struct {
	net.Conn
	isClosed bool
}

func (cc *CloseConn) Close() error {
	if !cc.isClosed {
		cc.isClosed = true
		return cc.Conn.Close()
	}

	return nil
}

type CloseListener struct {
	net.Listener
	isClosed bool
}

func (cl *CloseListener) Close() error {
	if !cl.isClosed {
		cl.isClosed = true
		return cl.Listener.Close()
	}
	return nil
}

func NewStreamServer() *StreamServer {
	tcpport := port.TcpPort()

	addr := ":" + strconv.Itoa(tcpport)

	cfg := config.GetCBtl()
	cfg.SetStreamPort(tcpport)

	ss := &StreamServer{addr: addr}
	ss.quit = make(chan struct{})

	return ss
}

func (ss *StreamServer) StartServer() error {
	lis, err := net.Listen("tcp", ss.addr)
	if err != nil {
		panic("failed to listen on %s" + ss.addr + " : " + err.Error())
	}

	ss.lis = &CloseListener{Listener: lis}
	defer ss.lis.Close()

	log.Println("Stream Server start at ", ss.addr)

	ss.wg.Add(1)
	go ss.serve()

	ss.wg.Wait()

	return nil
}

func (ss *StreamServer) serve() {
	defer ss.wg.Done()

	for {
		conn, err := ss.lis.Accept()
		if err != nil {
			select {
			case <-ss.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			ss.wg.Add(1)
			go func() {
				cc := &CloseConn{Conn: conn}
				ss.handleConn(cc)
			}()
		}

	}

}

func (ss *StreamServer) handleConn(conn net.Conn) {
	defer ss.wg.Done()
	defer conn.Close()

	conn.(*CloseConn).Conn.(*net.TCPConn).SetKeepAlive(true)

	if cs, err := handshake(conn); err != nil {
		return
	} else {
		var tgt Addr
		tgt, err = readAddr(cs)
		if err != nil {
			return
		}
		var rc net.Conn

		rc, err := net.Dial("tcp", tgt.String())
		if err != nil {
			return
		}
		defer rc.Close()
		rc.(*net.TCPConn).SetKeepAlive(true)
		log.Println("proxy ", cs.RemoteAddr().String(), "<->", tgt.String())
		err = relay2(cs, rc)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				return // ignore i/o timeout
			}
			log.Println("relay error: ", err)
		}
	}
}

func relay2(left, right net.Conn) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer func() {
			right.SetDeadline(time.Now())
			left.SetDeadline(time.Now())
			wg.Done()
		}()
		for {
			buf := stream.NewStreamBuf()
			n, err := left.Read(buf)
			if err != nil {
				return
			}
			var nw int
			nw, err = right.Write(buf[:n])
			if n != nw || err != nil {
				return
			}
		}
	}()
	for {
		buf := stream.NewStreamBuf()
		n, err := right.Read(buf)
		if err != nil {
			return err
		}
		var nw int
		nw, err = left.Write(buf[:n])
		if n != nw || err != nil {
			return err
		}
	}

	return nil
}

//func relay(left, right net.Conn) (int64, int64, error) {
//	type res struct {
//		N   int64
//		Err error
//	}
//	ch := make(chan res)
//
//	go func() {
//		n, err := io.Copy(right, left)
//		right.SetDeadline(time.Now())
//		left.SetDeadline(time.Now())
//		ch <- res{n, err}
//	}()
//
//	n, err := io.Copy(left, right)
//	right.SetDeadline(time.Now())
//	left.SetDeadline(time.Now())
//	rs := <-ch
//
//	if err == nil {
//		err = rs.Err
//	}
//	return n, rs.N, err
//}

func (ss *StreamServer) StopServer() {
	close(ss.quit)
	ss.lis.Close()

}
