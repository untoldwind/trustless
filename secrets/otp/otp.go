package otp

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"hash"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type OTP interface {
	// Authenticate verifies the OTP userCode.
	Authenticate(userCode string) bool
	// GetUserCode returns the current OTP userCode.
	GetUserCode() (string, time.Duration)
	// MaxDuration return the maximum duration a userCode might be valid
	MaxDuration() time.Duration
	GetEncodedSecret() string
	SetEncodedSecret(secret string) error
	// GetURL returns otpauth url
	GetURL() *url.URL
}

func ParseURL(urlStr string) (OTP, error) {
	otpauthURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid url")
	}
	if otpauthURL.Scheme != "otpauth" {
		return nil, errors.Errorf("Invalid otpauth schema: %s", otpauthURL.Scheme)
	}
	values := otpauthURL.Query()
	if len(values["secret"]) == 0 {
		return nil, errors.Errorf("No secret")
	}
	secret, err := base32.StdEncoding.DecodeString(strings.ToUpper(values.Get("secret")))
	if err != nil {
		return nil, errors.Wrap(err, "Invalid secret")
	}
	var label string
	if len(otpauthURL.Path) > 0 {
		label = otpauthURL.Path[1:]
	}
	issuer := values.Get("issuer")
	var hash func() hash.Hash
	switch values.Get("algorithm") {
	case "":
		hash = sha1.New
	case "SHA1":
		hash = sha1.New
	case "SHA256":
		hash = sha256.New
	case "SHA512":
		hash = sha512.New
	default:
		return nil, errors.Errorf("Invalid algorithm: %s", values.Get("algorithm"))
	}
	digits := 6
	if digitsParam := values.Get("digits"); len(digitsParam) > 0 {
		digits, err = strconv.Atoi(digitsParam)
		if err != nil {
			return nil, errors.Wrap(err, "Invalid digits")
		}
	}

	switch otpauthURL.Host {
	case "totp":
		totp := NewTOTP(secret)
		totp.Label = label
		totp.Issuer = issuer
		totp.Hash = hash
		totp.Digits = uint8(digits)

		if period := values.Get("period"); len(period) > 0 {
			seconds, err := strconv.Atoi(period)
			if err != nil {
				return nil, errors.Wrap(err, "Invalid period")
			}
			totp.TimeStep = time.Duration(seconds) * time.Second
		}
		return totp, nil
	default:
		return nil, errors.Errorf("Unsupported otp type: %s", otpauthURL.Host)
	}
}
