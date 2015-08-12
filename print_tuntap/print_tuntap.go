// Simple use of the tuntap package that prints packets received by the interface.
package main

import (
	"fmt"
	"github.com/liudanking/tuntap"
	"github.com/liudanking/tuntap/util"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("syntax:", os.Args[0], "tun|tap", "<device name>")
		return
	}

	var typ tuntap.DevKind
	switch os.Args[1] {
	case "tun":
		typ = tuntap.DevTun
	case "tap":
		typ = tuntap.DevTap
	default:
		fmt.Println("Unknown device type", os.Args[1])
		return
	}

	tun, err := tuntap.Open(os.Args[2], typ)
	if err != nil {
		fmt.Println("Error opening tun/tap device:", err)
		return
	}

	fmt.Println("Listening on", tun.Name())
	fmt.Println([]byte(tun.Name()))
	buf := make([]byte, 1522)
	for {
		// pkt, err := tun.ReadPacket()
		_, err := tun.Read(buf)

		if err != nil {
			fmt.Println("Read error:", err)
		} else {
			fmt.Printf("from %s to %s\n", util.IPv4Source(buf).String(), util.IPv4Destination(buf).String())
			// if pkt.Truncated {
			// 	fmt.Printf("!")
			// } else {
			// 	fmt.Printf(" ")
			// }
			// fmt.Printf("%x %x\n", pkt.Protocol, pkt.Packet)

		}
	}
}
