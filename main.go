package main

import (
	"net"

	"github.com/bzp2010/knockdoor/knockdoor"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

var (
	visitors = map[string]*knockdoor.Visitor{}
)

func main() {
	ip, _ := net.ResolveIPAddr("ip4", "0.0.0.0")
	conn, err := net.ListenIP("ip4:tcp", ip)

	if err != nil {
		panic(err)
	}

	ipRawConn, _ := ipv4.NewRawConn(conn)
	for {
		buf := make([]byte, 1500)
		hdr, payload, _, _ := ipRawConn.ReadFrom(buf)
		sourceIP := hdr.Src.String()

		// skip localhost
		if sourceIP == "127.0.0.1" {
			continue
		}

		packet := gopacket.NewPacket(payload, layers.LayerTypeTCP, gopacket.Default)
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)

			if !tcp.SYN {
				continue
			}

			if _, ok := visitors[sourceIP]; !ok {
				visitors[sourceIP] = knockdoor.NewVisitor()
			}
			visitors[sourceIP].Handle(tcp)
		}
	}
}

/* fmt.Println("SOURCE:", hdr.Src, tcp.SrcPort)
fmt.Println("DEST:", hdr.Dst, tcp.DstPort)
fmt.Println("TCP SYN:", tcp.SYN)
fmt.Println("TCP ACK:", tcp.ACK)
fmt.Println("TCP FIN:", tcp.FIN)
fmt.Println("TCP RST:", tcp.RST)
fmt.Println("TCP PSH:", tcp.PSH)
fmt.Println("TCP URG:", tcp.URG)
fmt.Println("TCP ECE:", tcp.ECE)
fmt.Println("TCP CWR:", tcp.CWR)
fmt.Println("TCP NS:", tcp.NS)
fmt.Println("TCP Seq:", tcp.Seq) */
