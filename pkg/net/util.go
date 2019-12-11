package net

func Cksum16(b []byte, n int, init uint32) uint16 {
	sum := init
	for i := 0; i < n-1; i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
		if (sum >> 16) > 0 {
			sum = (sum & 0xffff) + (sum >> 16)
		}
	}
	if n&1 != 0 {
		sum += uint32(b[n-1]) << 8
		if (sum >> 16) > 0 {
			sum = (sum & 0xffff) + (sum >> 16)
		}
	}
	return ^(uint16(sum))
}
