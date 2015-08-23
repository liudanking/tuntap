package tuntap

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	flagTruncated = 0x1
	// flags contains the flags that tell the kernel which kind of interface we want (tun or tap).
	// Basically, it can either take the value IFF_TUN to indicate a TUN device (no ethernet headers
	// in the packets), or IFF_TAP to indicate a TAP device (with ethernet headers in packets).
	// Additionally, another flag IFF_NO_PI can be ORed with the base value.
	// IFF_NO_PI tells the kernel to not provide packet information.
	// The purpose of IFF_NO_PI is to tell the kernel that packets will be "pure" IP packets,
	// with no added bytes. Otherwise (if IFF_NO_PI is unset), 4 extra bytes are added to the
	// beginning of the packet (2 flag bytes and 2 protocol bytes).
	// IFF_NO_PI need not match between interface creation and reconnection time.
	// Also note that when capturing traffic on the interface with Wireshark, those 4 bytes are never shown.
	iffTun  = 0x1
	iffTap  = 0x2
	iffNoPi = 0x1000
	// https://www.mail-archive.com/user-mode-linux-devel@lists.sourceforge.net/msg00475.html
	iffOneQueue = 0x2000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func Open(ifPattern string, kind DevKind) (*Interface, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	ifName, err := createInterface(file, ifPattern, kind)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &Interface{kind, ifName, file}, nil
}

func createInterface(file *os.File, ifPattern string, kind DevKind) (string, error) {
	var req ifReq
	req.Flags = iffNoPi
	copy(req.Name[:], ifPattern)
	switch kind {
	case DevTun:
		req.Flags |= iffTun
	case DevTap:
		req.Flags |= iffTap
	default:
		panic("Unknown interface type")
	}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if err != 0 {
		return "", err
	}
	return strings.TrimRight(string(req.Name[:4]), "\x00"), nil
}
