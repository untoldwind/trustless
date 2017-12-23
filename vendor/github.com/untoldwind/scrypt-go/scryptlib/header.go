package scryptlib

type Header struct {
	Magic      [6]byte
	Version    uint8
	Params     Params
	Salt       [32]byte
	HeaderHash [16]byte
	HeaderHMAC [32]byte
}

var Magic = [6]byte{'s', 'c', 'r', 'y', 'p', 't'}
