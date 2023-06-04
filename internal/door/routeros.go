package door

import (
	"fmt"

	"github.com/bzp2010/knockdoor/internal/config"
	"github.com/bzp2010/knockdoor/internal/log"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
)

type routerOSDoor struct {
	cfg config.Door
}

// NewRouterOSDoor creates a new RouterOS door (address list)
func NewRouterOSDoor(cfg config.Door) Door {
	return &routerOSDoor{
		cfg: cfg,
	}
}

// Open opens the door
func (r *routerOSDoor) Open(ip string) error {
	log.GetLogger().Infow("Add visitor IP to RouterOS address list", "ip", ip, "addressList", r.cfg.RouterOS.AddressListName)
	request := req.
		DefaultClient().
		SetCommonBasicAuth(r.cfg.RouterOS.Username, r.cfg.RouterOS.Password).R()
	resp := request.
		SetBodyJsonMarshal(map[string]string{"list": r.cfg.RouterOS.AddressListName, "address": ip, "timeout": "1d"}).
		MustPut(fmt.Sprintf("%s%s", r.cfg.RouterOS.Endpoint, "/rest/ip/firewall/address-list"))

	if resp.StatusCode >= 300 {
		return errors.Errorf("failed to open the door: %s", resp.String())
	}
	log.GetLogger().Infow("Visitor IP added to RouterOS address list", "ip", ip, "addressList", r.cfg.RouterOS.AddressListName)
	return nil
}
