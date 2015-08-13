// Simple use of the tuntap package that prints packets received by the interface.
package main

import (
	"encoding/binary"
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

	fmt.Println("Listening on", string(tun.Name()))
	buf := make([]byte, 1522)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			fmt.Println("Read error:", err)
		} else {
			if util.IsIPv4(buf) {
				fmt.Printf("%d bytes from iface, IHL:%02X, TTL:%d\n", n, buf[0], buf[8])
				fmt.Printf("from %s to %s\n", util.IPv4Source(buf).String(), util.IPv4Destination(buf).String())
				fmt.Printf("protocol %02x checksum %02x\n", util.IPv4Protocol(buf), binary.BigEndian.Uint16(buf[22:24]))
				if util.IPv4Protocol(buf) == util.ICMP {
					srcip := make([]byte, 4)
					copy(srcip, buf[12:16])
					copy(buf[12:16], buf[16:20])
					copy(buf[16:20], srcip)
					buf[20] = 0x00
					buf[21] = 0x00

					buf[22] = 0x00
					buf[23] = 0x00

					var checksum uint32 = 0
					for i := 20; i < n; i += 2 {
						checksum += uint32(binary.BigEndian.Uint16(buf[i : i+2]))
					}

					checksum = (checksum >> 16) + (checksum & 0xffff)
					checksum += (checksum >> 16)
					checksum = 0xffff &^ checksum
					fmt.Printf("my checksum:%02x\n", uint16(checksum))
					buf[22] = byte((checksum & 0xff00) >> 8)
					buf[23] = byte(checksum & 0xff)
					fmt.Printf("rsp: from %s to %s\n", util.IPv4Source(buf).String(), util.IPv4Destination(buf).String())

					_, err = tun.Write(buf)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}
