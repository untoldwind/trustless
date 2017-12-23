package scryptlib

import (
	"time"

	"golang.org/x/crypto/scrypt"
)

// cpuperf estimates how many salsa20/8 core operations can be done per second
func cpuperf() (int64, error) {
	var diff time.Duration = 0
	var i int64 = 0
	start := time.Now()

	for {
		if _, err := scrypt.Key(nil, nil, 128, 1, 1, 0); err != nil {
			return 0, err
		}
		i += 512
		diff = time.Since(start)
		if diff > 10*time.Millisecond {
			break
		}
	}

	return i * int64(time.Second) / int64(diff), nil
}

func pickparams(memlimit uint64, maxtime float64) (*Params, error) {
	ops, err := cpuperf()
	if err != nil {
		return nil, err
	}
	opslimit := float64(ops) * maxtime
	if opslimit < 32768 {
		opslimit = 32768
	}
	R := uint32(8)
	if opslimit < float64(memlimit/32) {
		maxN := int(opslimit / float64(R*4))
		logN := uint8(1)
		for ; logN < 63; logN++ {
			if (2 << logN) > maxN {
				break
			}
		}
		return &Params{
			LogN: logN,
			R:    R,
			P:    1,
		}, nil
	} else {
		maxN := int(memlimit / uint64(R*120))
		logN := uint8(1)
		for ; logN < 63; logN++ {
			if (2 << logN) > maxN {
				break
			}
		}

		maxrp := uint64(opslimit/4) / uint64(1<<logN)
		if maxrp > 0x3fffffff {
			maxrp = 0x3fffffff
		}
		return &Params{
			LogN: logN,
			R:    R,
			P:    uint32(maxrp),
		}, nil
	}
}
