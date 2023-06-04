package main

import (
	"fmt"
	"os"

	"github.com/bzp2010/knockdoor/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
