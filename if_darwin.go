// Requirement: Install driver from http://tuntaposx.sourceforge.net/
package tuntap

import (
	"os"
)

const (
	flagTruncated = 0x1
	iffTun        = 0x1
	iffTap        = 0x2
	iffNoPi       = 0x1000
	iffOneQueue   = 0x2000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func Open(ifPattern string, kind DevKind) (*Interface, error) {
	file, err := os.OpenFile("/dev/"+ifPattern, os.O_RDWR, os.ModeCharDevice)
	if err != nil {
		return nil, err
	}

	return &Interface{kind, ifPattern, file}, nil
}
