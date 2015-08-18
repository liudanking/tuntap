// [15-08-18 23:22:43][INFO] ip -4 route add 0.0.0.0/1 via 10.0.0.40 dev tun0
// [15-08-18 23:22:43][INFO] ip -4 route add 128.0.0.0/1 via 10.0.0.40 dev tun0
// make all traffic via tun0
package main

import (
	"fmt"
	"github.com/liudanking/tuntap"
	"github.com/liudanking/tuntap/util"
	"net"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	tun, err := tuntap.Open("tun0", tuntap.DevTun)
	checkError(err)
	laddr, err := net.ResolveUDPAddr("udp", "192.168.102.32:0")
	checkError(err)
	raddr, err := net.ResolveUDPAddr("udp", "192.168.102.31:9826")
	checkError(err)

	conn, err := net.DialUDP("udp", laddr, raddr)
	checkError(err)
	defer conn.Close()

	fmt.Println("Waiting IP Packet from tun interface")
	go func() {
		buf := make([]byte, 10000)
		for {
			n, err := tun.Read(buf)
			if err != nil {
				fmt.Println("tun Read error:", err)
				continue
			}
			fmt.Printf("receive %d bytes, from %s to %s, \n", n, util.IPv4Source(buf).String(), util.IPv4Destination(buf).String())
			n, err = conn.Write(buf[:n])
			if err != nil {
				fmt.Println("udp write error:", err)
				continue
			}
			fmt.Printf("write %d bytes to udp network\n", n)
		}
	}()

	buf := make([]byte, 10000)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("udp Read error:", err)
			continue
		}
		fmt.Sprintf("receive %d bytes, from %s to %s, \n", n, util.IPv4Source(buf).String(), util.IPv4Destination(buf).String())
		n, err = tun.Write(buf[:n])
		if err != nil {
			fmt.Println("udp write error:", err)
			continue
		}
		fmt.Printf("write %d bytes to tun interface\n", n)
	}
}
