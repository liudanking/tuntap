// Todo: support Mac OS X
package tuntap

import (
	"os"
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
	return "", nil
}
