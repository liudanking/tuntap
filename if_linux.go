package tuntap

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	flagTruncated = 0x1

	iffTun      = 0x1
	iffTap      = 0x2
	iffOneQueue = 0x2000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func createInterface(file *os.File, ifPattern string, kind DevKind) (string, error) {
	var req ifReq
	req.Flags = iffOneQueue
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
	return strings.TrimRight(string(req.Name[:]), "\x00"), nil
}
