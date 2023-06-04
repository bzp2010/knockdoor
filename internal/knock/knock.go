package knock

import (
	"strconv"

	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

// Knock is an interface used to represent the "knock" mechanism
type Knock interface {
	Handle(ipHeader *ipv4.Header, tcpPacket *layers.TCP) bool
}

func portToString(port layers.TCPPort) string {
	return strconv.Itoa(int(port))
}
