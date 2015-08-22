package main

import (
	"fmt"
	"github.com/liudanking/tuntap"
	"net"
	"os"
	"os/exec"
	"runtime"
)

const (
	SERVER_IP = "10.0.0.2"
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
	exeCmd(fmt.Sprintf("ip addr add %s/24 dev tun0", SERVER_IP))
	exeCmd("iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -o eth0 -j MASQUERADE")
}

func setTunDarwin() {
	exeCmd(fmt.Sprintf("ifconfig tun0 inet %s/24 %s up", SERVER_IP, SERVER_IP))
	exeCmd(fmt.Sprintf("route -n add 10.0.0.0/24 %s", SERVER_IP))
	exec.Command("bash", "-c", `echo "nat on en0 inet from 10.0.0.0/24 to any -> en0" | pfctl -v -ef -`).Output()
}

func main() {
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
