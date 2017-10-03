package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"hash"
	"net/url"
	"strconv"
	"time"
)

// TOTP contains the different configurable values for a given TOTP
// invocation.
type TOTP struct {
	Time      func() time.Time
	Tries     []int64
	TimeStep  time.Duration
	Digits    uint8
	Hash      func() hash.Hash
	Label     string
	Issuer    string
	secretKey []byte
}

// NewTOTP creates a time based OTP from a secret with default options
func NewTOTP(secretKey []byte) *TOTP {
	return &TOTP{
		Time:      time.Now,
		Tries:     []int64{0, -1},
		TimeStep:  30 * time.Second,
		Digits:    6,
		Hash:      sha1.New,
		secretKey: secretKey,
	}
}

var digitModulo = []int64{
	1,          // 0
	10,         // 1
	100,        // 2
	1000,       // 3
	10000,      // 4
	100000,     // 5
	1000000,    // 6
	10000000,   // 7
	100000000,  // 8
	1000000000, // 9
}

func (o *TOTP) Authenticate(userCode string) bool {

	if int(o.Digits) != len(userCode) {
		return false
	}

	uc, err := strconv.ParseInt(userCode, 10, 64)
	if err != nil {
		return false
	}

	t := o.Time().Unix() / int64(o.TimeStep/time.Second)

	for i := 0; i < len(o.Tries); i++ {
		b := t + o.Tries[i]

		code := o.calculateCode(b)

		if code == uc {
			return true
		}
	}

	return false
}

func (o *TOTP) GetUserCode() (string, error) {
	t := o.Time().Unix() / int64(o.TimeStep/time.Second)

	code := o.calculateCode(t)

	return fmt.Sprintf("%06d", code), nil
}

func (o *TOTP) GetURL() (*url.URL, error) {
	return nil, nil
}

func (o *TOTP) calculateCode(time int64) int64 {
	var tbuf [8]byte

	hm := hmac.New(o.Hash, o.secretKey)
	var hashbuf []byte
	tbuf[0] = byte(time >> 56)
	tbuf[1] = byte(time >> 48)
	tbuf[2] = byte(time >> 40)
	tbuf[3] = byte(time >> 32)
	tbuf[4] = byte(time >> 24)
	tbuf[5] = byte(time >> 16)
	tbuf[6] = byte(time >> 8)
	tbuf[7] = byte(time)

	hm.Reset()
	hm.Write(tbuf[:])
	hashbuf = hm.Sum(hashbuf[:0])

	offset := hashbuf[len(hashbuf)-1] & 0xf
	truncatedHash := hashbuf[offset:]

	code := int64(truncatedHash[0])<<24 |
		int64(truncatedHash[1])<<16 |
		int64(truncatedHash[2])<<8 |
		int64(truncatedHash[3])

	code &= 0x7FFFFFFF
	code %= digitModulo[int(o.Digits)]

	return code
}