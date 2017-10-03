package otp_test

import (
	"testing"

	"github.com/untoldwind/trustless/secrets/otp"
)

func TestGetUserCode(t *testing.T) {
	otp := otp.NewTOTP([]byte("youCanTrustMe"))
	code, _ := otp.GetUserCode()
	auth := otp.Authenticate(code)
	if !auth {
		t.Errorf("GetUserCode didn't authenticate correctly.")
	}
}
