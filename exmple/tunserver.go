// ip link set dev tun0 up
// ip addr add 10.0.0.30/24 dev tun0
// iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -o eth0 -j MASQUERADE
package main

import (
	"fmt"
	"github.com/liudanking/tuntap"
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
	addr, err := net.ResolveUDPAddr("udp", ":9826")
	checkError(err)
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	defer conn.Close()
	raddr := &net.UDPAddr{}
	fmt.Println("Waiting IP Packet from UDP")
	go func() {
		buf := make([]byte, 10000)
		for {
			n, fromAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("ReadFromUDP error:", err)
				continue
			}
			raddr = fromAddr
			fmt.Printf("receive %d bytes from %s\n", n, fromAddr.String())
			n, _ = tun.Write(buf[:n])
			fmt.Printf("write %d bytes to tun interface\n", n)
		}
	}()

	buf := make([]byte, 10000)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			fmt.Println("run read error:", err)
			continue
		}
		n, err = conn.WriteTo(buf[:n], raddr)
		// n, err = conn.Write(buf[:n])
		fmt.Printf("write %d bytes to udp network\n", n)
	}
}
