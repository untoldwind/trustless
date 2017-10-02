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
	GetUserCode() (string, error)
	// GetURL returns otpauth url
	GetURL() (*url.URL, error)
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
	secret, err := base32.StdEncoding.DecodeString(strings.ToUpper(values["secret"][0]))
	if err != nil {
		return nil, errors.Wrap(err, "Invalid secret")
	}
	var label string
	if len(otpauthURL.Path) > 0 {
		label = otpauthURL.Path[1:]
	}
	var issuer string
	if issuers := values["issuer"]; len(issuers) > 0 {
		issuer = issuers[0]
	}
	var hash func() hash.Hash = sha1.New
	if algorithm := values["algorithm"]; len(algorithm) > 0 {
		switch algorithm[0] {
		case "SHA1":
		case "SHA256":
			hash = sha256.New
		case "SHA512":
			hash = sha512.New
		default:
			return nil, errors.Errorf("Invalid algorithm: %s", algorithm[0])
		}
	}
	digits := 6
	if digitsParam := values["digits"]; len(digitsParam) > 0 {
		digits, err = strconv.Atoi(digitsParam[0])
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

		if period := values["period"]; len(period) > 0 {
			seconds, err := strconv.Atoi(period[0])
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
