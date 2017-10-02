package otp_test

import (
	"testing"

	"github.com/untoldwind/trustless/secrets/otp"
)

func TestGetUserCode(t *testing.T) {
	otp := otp.NewTOTP([]byte("youCanTrustMe"))
	code, err := otp.GetUserCode()
	if err != nil {
		t.Errorf("GetUserCode failed: %s", err.Error())
	}
	auth := otp.Authenticate(code)
	if !auth {
		t.Errorf("GetUserCode didn't authenticate correctly.")
	}
}
