package knock

import (
	"github.com/bzp2010/knockdoor/internal/config"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
)

// Visitor is a visitor that want to establish a real TCP connection
type Visitor struct {
	mechanism Knock
	sourceIP  string
}

// NewVisitor creates a new visitor to knock the door
func NewVisitor(cfg config.Port, doneCallback func(ip string)) *Visitor {
	v := &Visitor{}

	switch cfg.Mode {
	case "totp":
		v.mechanism = NewTOTPKnock(*cfg.TOTP, func() {
			doneCallback(v.sourceIP)
		})
	case "static":
		fallthrough
	default:
		v.mechanism = NewStaticKnock(*cfg.Static, func() {
			doneCallback(v.sourceIP)
		})
	}
	return v
}

// Handle handles the TCP packet
func (v *Visitor) Handle(ipHeader *ipv4.Header, tcpPacket *layers.TCP) bool {
	v.sourceIP = ipHeader.Src.String()
	return v.mechanism.Handle(ipHeader, tcpPacket)
}
