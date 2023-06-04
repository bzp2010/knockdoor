package knockdoor

import (
	"strconv"

	"github.com/google/gopacket/layers"
	"github.com/looplab/fsm"
)

var (
	portSerial  = []layers.TCPPort{9999, 9998, 9997}
	initialPort = portSerial[0]

	fsmEvents = fsm.Events{}
)

func init() {
	for i, port := range portSerial {
		if i == 0 {
			fsmEvents = append(fsmEvents, fsm.EventDesc{Name: "knock", Src: []string{"new_visitor"}, Dst: "STAGE_" + portToString(port)})
		} else if i == len(portSerial)-1 {
			fsmEvents = append(fsmEvents, fsm.EventDesc{Name: "knock_" + portToString(port), Src: []string{"STAGE_" + portToString(portSerial[i-1])}, Dst: "STAGE_OPEN_DOOR"})
		} else {
			fsmEvents = append(fsmEvents, fsm.EventDesc{Name: "knock_" + portToString(port), Src: []string{"STAGE_" + portToString(portSerial[i-1])}, Dst: "STAGE_" + portToString(port)})
		}
	}
}

func portToString(port layers.TCPPort) string {
	return strconv.Itoa(int(port))
}
