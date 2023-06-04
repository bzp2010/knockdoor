package knockdoor

import (
	"context"
	"fmt"

	"github.com/google/gopacket/layers"
	"github.com/looplab/fsm"
	"golang.org/x/net/ipv4"
)

// Visitor is a visitor that want to establish a real TCP connection
type Visitor struct {
	new  bool
	done bool
	fsm  *fsm.FSM

	sourceIP string
}

// NewVisitor creates a new visitor to knock the door
func NewVisitor() *Visitor {
	v := &Visitor{
		new: true,
	}
	v.fsm = v.generateFSM(portSerial)
	return v
}

// Handle handles the tcp packet
func (v *Visitor) Handle(ipHeader *ipv4.Header, tcpPacket *layers.TCP) bool {
	v.sourceIP = ipHeader.Src.String()
	if v.new {
		if tcpPacket.DstPort != initialPort {
			return true
		}
		v.new = false
		v.fsm.Event(context.Background(), "knock")
	} else {
		v.fsm.Event(context.Background(), "knock_"+portToString(tcpPacket.DstPort))
	}
	return v.done
}

func (v *Visitor) generateFSM(ports []layers.TCPPort) *fsm.FSM {
	return fsm.NewFSM("new_visitor", fsmEvents, fsm.Callbacks{
		"enter_state": func(_ context.Context, e *fsm.Event) {
			if e.Dst == "STAGE_OPEN_DOOR" {
				v.done = true
				openDoor(v.sourceIP)
			}
		},
	})
}

func openDoor(ip string) {
	fmt.Printf("TODO: SRC_IP: %s, open the door\n", ip)
}
