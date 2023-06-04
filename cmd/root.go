package cmd

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/bzp2010/knockdoor/internal/config"
	"github.com/bzp2010/knockdoor/internal/knock"
	"github.com/bzp2010/knockdoor/internal/log"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/ipv4"
)

var (
	configFile = "config/config.yaml"
)

var (
	visitorsMutex sync.Mutex
	visitors      = map[string]*knock.Visitor{}
)

// NewRootCommand for main package
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Timed task scheduler",
		RunE:  run,
	}

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "config file")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	// setup config
	cfg := config.NewDefaultConfig()
	if err := config.SetupConfig(&cfg, configFile); err != nil {
		return errors.Wrap(err, "failed to setup config")
	}

	// setup logger
	if err := log.SetupLogger(cfg); err != nil {
		return errors.Wrap(err, "failed to setup logger")
	}

	err := listenIPv4(cfg)
	if err != nil {
		return err
	}

	go visitorCleaner()

	// block main goroutine
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	return nil
}

func listenIPv4(cfg config.Config) error {
	ip, _ := net.ResolveIPAddr("ip4", "0.0.0.0")
	conn, err := net.ListenIP("ip4:tcp", ip)
	if err != nil {
		return errors.Wrap(err, "failed to listen on raw socket")
	}
	ipRawConn, _ := ipv4.NewRawConn(conn)
	go handleIPv4RawConn(cfg, ipRawConn)

	return nil
}

func handleIPv4RawConn(cfg config.Config, conn *ipv4.RawConn) {
	for {
		buf := make([]byte, 1500)
		ipHdr, payload, _, _ := conn.ReadFrom(buf)
		sourceIP := ipHdr.Src.String()

		// skip loopback
		if cfg.Port.SkipLoopback && ipHdr.Src.IsLoopback() {
			continue
		}

		packet := gopacket.NewPacket(payload, layers.LayerTypeTCP, gopacket.Default)
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)

			if !tcp.SYN {
				continue
			}

			if _, ok := visitors[sourceIP]; !ok {
				visitorsMutex.Lock()
				visitors[sourceIP] = knock.NewVisitor(cfg.Port, func(ip string) {
					fmt.Println("OPEN THE DOOR, I AM", ip)
				})
				visitorsMutex.Unlock()
			}

			if clean := visitors[sourceIP].Handle(ipHdr, tcp); clean {
				visitorsMutex.Lock()
				delete(visitors, sourceIP)
				visitorsMutex.Unlock()
			}
		}
	}
}

func visitorCleaner() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		<-ticker.C
		visitorsMutex.Lock()
		visitors = map[string]*knock.Visitor{}
		visitorsMutex.Unlock()
	}
}
