package knock

import (
	"context"

	"github.com/bzp2010/knockdoor/internal/log"
	"github.com/google/gopacket/layers"
	"github.com/looplab/fsm"
	"golang.org/x/net/ipv4"
)

type staticKnock struct {
	portSerial []layers.TCPPort

	fsm  *fsm.FSM
	new  bool
	done bool
}

// NewStaticKnock creates a new static knock
func NewStaticKnock(ports []uint16, doneCallback func()) Knock {
	// convert TCP port
	portSerial := make([]layers.TCPPort, len(ports))
	for i, port := range ports {
		portSerial[i] = layers.TCPPort(port)
	}

	// create events
	fsmEvents := fsm.Events{}
	for i, port := range portSerial {
		if i == 0 {
			fsmEvents = append(fsmEvents, fsm.EventDesc{Name: "knock", Src: []string{"new_visitor"}, Dst: "STAGE_" + portToString(port)})
		} else if i == len(portSerial)-1 {
			fsmEvents = append(fsmEvents, fsm.EventDesc{Name: "knock_" + portToString(port), Src: []string{"STAGE_" + portToString(portSerial[i-1])}, Dst: "STAGE_OPEN_DOOR"})
		} else {
			fsmEvents = append(fsmEvents, fsm.EventDesc{Name: "knock_" + portToString(port), Src: []string{"STAGE_" + portToString(portSerial[i-1])}, Dst: "STAGE_" + portToString(port)})
		}
	}

	knock := &staticKnock{
		portSerial: portSerial,
		new:        true,
	}
	knock.fsm = knock.generateFSM(fsmEvents, doneCallback)
	return knock
}

func (s *staticKnock) Handle(ipHeader *ipv4.Header, tcpPacket *layers.TCP) bool {
	if s.new {
		if tcpPacket.DstPort != s.portSerial[0] {
			log.GetLogger().Debugw("New visitor arrivals, but the port does not match the starting port",
				"ip", ipHeader.Src.String(),
				"port", portToString(tcpPacket.DstPort),
			)
			return true
		}
		log.GetLogger().Infow("New visitor arrivals", "ip", ipHeader.Src.String())
		s.new = false
		s.fsm.Event(context.Background(), "knock")
	} else {
		s.fsm.Event(context.Background(), "knock_"+portToString(tcpPacket.DstPort))
	}
	return s.done
}

func (s *staticKnock) generateFSM(fsmEvents fsm.Events, doneCallback func()) *fsm.FSM {
	return fsm.NewFSM("new_visitor", fsmEvents, fsm.Callbacks{
		"enter_state": func(_ context.Context, e *fsm.Event) {
			if e.Dst == "STAGE_OPEN_DOOR" {
				s.done = true
				doneCallback()
			}
		},
	})
}
