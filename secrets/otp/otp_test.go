package otp_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/secrets/otp"
)

func TestParse(t *testing.T) {
	require := require.New(t)

	parsed, err := otp.ParseURL("otpauth://totp/ACME%20Co:john@example.com?secret=MBUO5TLQRRQU4ZFZZE4Q47NY5RAKUXLN&issuer=ACME%20Co&algorithm=SHA1&digits=6&period=30")
	require.Nil(err)

	topt, ok := parsed.(*otp.TOTP)
	require.True(ok)

	require.Equal("ACME Co", topt.Issuer)
	require.Equal("ACME Co:john@example.com", topt.Label)

	topt.Time = func() time.Time {
		return time.Unix(1507030530, 0)
	}

	code, _ := topt.GetUserCode()
	require.Equal("676940", code)

	genUrl := topt.GetURL()
	require.NotNil(genUrl)
	require.Equal("otpauth://totp/ACME%20Co:john@example.com?algorithm=SHA1&digits=6&issuer=ACME+Co&period=30&secret=MBUO5TLQRRQU4ZFZZE4Q47NY5RAKUXLN", genUrl.String())
}
