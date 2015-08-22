package main

import (
	"flag"
	"fmt"
	"github.com/liudanking/tuntap"
	"github.com/liudanking/tuntap/util"
	"net"
	"os"
	"os/exec"
	"runtime"
)

const (
	CLIENT_IP = "10.0.0.30"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func exeCmd(cmd string) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("execute %s error:%v", cmd, err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

func setTunLinux() {
	exeCmd("ip link set dev tun0 up")
	exeCmd(fmt.Sprintf("ip addr add %s/24 dev tun0", CLIENT_IP))
	// make all traffic via tun0
	exeCmd(fmt.Sprintf("ip -4 route add 0.0.0.0/1 via %s dev tun0", CLIENT_IP))
	exeCmd(fmt.Sprintf("ip -4 route add 128.0.0.0/1 via %s dev tun0", CLIENT_IP))
}

func setTunDarwin() {
	exeCmd(fmt.Sprintf("ifconfig tun0 inet %s/24 %s up", CLIENT_IP, CLIENT_IP))
	exeCmd(fmt.Sprintf("route -n add 0.0.0.0/1 %s", CLIENT_IP))
	exeCmd(fmt.Sprintf("route -n add 128.0.0.0/1 %s", CLIENT_IP))
}

func main() {
	lip := flag.String("l", "192.168.102.32", "client local ip")
	rip := flag.String("r", "192.168.102.31", "server remote ip")
	flag.Parse()

	tun, err := tuntap.Open("tun0", tuntap.DevTun)
	checkError(err)
	switch runtime.GOOS {
	case "linux":
		setTunLinux()
	case "darwin":
		setTunDarwin()
	default:
		fmt.Println("OS NOT supported")
		os.Exit(1)
	}

	laddr, err := net.ResolveUDPAddr("udp", *lip+":0")
	checkError(err)
	raddr, err := net.ResolveUDPAddr("udp", *rip+":9826")
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
