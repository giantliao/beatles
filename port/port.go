package port

import (
	"encoding/binary"
	"github.com/giantliao/beatles/wallet"
	"github.com/kprc/libeth/account"
	"sync"
)

var (
	portBasic int
	portLock  sync.Mutex
)

const (
	HTTP_Server_Port int = iota
	TCP_Server_port
	UDP_Server_port
)

const (
	int32size int    = 4
	PortRange uint32 = 4096
	PortBasic int    = 44000
)

func TcpPort() int {
	if portBasic == 0 {
		portBasic = calcPortBasic()
	}
	return portBasic + TCP_Server_port
}

func HttpPort() int {
	if portBasic == 0 {
		portBasic = calcPortBasic()
	}

	return portBasic + HTTP_Server_Port
}

func UpdPort() int {
	if portBasic == 0 {
		portBasic = calcPortBasic()
	}

	return portBasic + UDP_Server_port
}

func calcPortBasic() int {
	w, err := wallet.GetWallet()
	if err != nil {
		panic("wallet not found")
	}

	var addr account.Address
	addr, _, _, err = w.BtlAddress().Address()
	if err != nil {
		panic("not a correct address")
	}

	cnt := len(addr) / int32size

	u := bytes2uint32(addr[:int32size])

	for i := 1; i < cnt; i++ {
		u = u ^ bytes2uint32(addr[int32size*i:int32size*(i+1)])
	}

	hu := (u >> 16 & 0xFFFF)
	lu := (u & 0xFFFF)

	u = hu ^ lu

	u = u & (PortRange - 1)

	return int(u) + PortBasic
}

func bytes2uint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}
