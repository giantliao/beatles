package streamserver

import (
	"encoding/hex"
	"fmt"
	"github.com/giantliao/beatles-protocol/stream"
	"net"
	"strconv"
)

const (
	AtypIPv4       = 1
	AtypDomainName = 3
	AtypIPv6       = 4
)

type Addr []byte

// String serializes SOCKS address a to string form.
func (a Addr) String() string {
	var host, port string

	switch a[0] { // address type
	case AtypDomainName:
		host = string(a[2 : 2+int(a[1])])
		port = strconv.Itoa((int(a[2+int(a[1])]) << 8) | int(a[2+int(a[1])+1]))
	case AtypIPv4:
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(a[1+net.IPv4len]) << 8) | int(a[1+net.IPv4len+1]))
	case AtypIPv6:
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(a[1+net.IPv6len]) << 8) | int(a[1+net.IPv6len+1]))
	}

	return net.JoinHostPort(host, port)
}

func readAddr(conn net.Conn) (Addr, error) {
	b := stream.NewStreamBuf()
	n, err := conn.Read(b)
	if err != nil {
		return nil, err
	}

	fmt.Println("target is ", Addr(b[:n]).String(), hex.EncodeToString(b[:n]))

	return b[:n], nil

}
