package util

import ()

func Checksum(data []byte) (sum uint16) {
	var cksum uint32 = 0
	length := len(data)
	for i := 0; i+1 < length; i += 2 {
		cksum += uint32(uint16(data[i])<<8 + uint16(data[i+1]))
	}
	cksum = (cksum >> 16) + (cksum & 0xffff)
	cksum += (cksum >> 16)
	cksum += (cksum >> 16)
	cksum = 0xffff &^ cksum
	return uint16(cksum)
}
