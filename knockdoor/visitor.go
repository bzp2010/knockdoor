package knockdoor

import (
	"context"
	"fmt"

	"github.com/google/gopacket/layers"
	"github.com/looplab/fsm"
)

// Visitor is a visitor that want to establish a real TCP connection
type Visitor struct {
	new bool
	fsm *fsm.FSM
}

// NewVisitor creates a new visitor to knock the door
func NewVisitor() *Visitor {
	return &Visitor{
		new: true,
		fsm: genFSM(portSerial),
	}
}

// Handle handles the tcp packet
func (v *Visitor) Handle(tcpPacket *layers.TCP) {
	if v.new {
		if tcpPacket.DstPort != initialPort {
			return
		}
		v.new = false
		v.fsm.Event(context.Background(), "knock")
	} else {
		v.fsm.Event(context.Background(), "knock_"+portToString(tcpPacket.DstPort))
	}
}

func genFSM(ports []layers.TCPPort) *fsm.FSM {
	return fsm.NewFSM("new_visitor", fsmEvents, fsm.Callbacks{
		"enter_state": func(_ context.Context, e *fsm.Event) {
			if e.Dst == "STAGE_OPEN_DOOR" {
				fmt.Println("TODO")
			}
		},
	})
}
