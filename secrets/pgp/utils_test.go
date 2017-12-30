package pgp

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesToWords(t *testing.T) {
	require := require.New(t)

	if sizeOfUInt != 4 && sizeOfUInt != 8 {
		t.FailNow()
	}

	raw := make([]byte, 8*sizeOfUInt)
	for i := 0; i < 8; i++ {
		raw[sizeOfUInt*i] = byte(i)
		raw[sizeOfUInt*i+1] = byte(i)
		raw[sizeOfUInt*i+2] = byte(i)
		raw[sizeOfUInt*i+3] = byte(i)
	}

	words := bytesToWords(raw)
	require.Len(words, 8)
	require.Equal(big.Word(0), words[0])
	require.Equal(big.Word(0x01010101), words[1])
	require.Equal(big.Word(0x02020202), words[2])
	require.Equal(big.Word(0x03030303), words[3])
	require.Equal(big.Word(0x04040404), words[4])
	require.Equal(big.Word(0x05050505), words[5])
	require.Equal(big.Word(0x06060606), words[6])
	require.Equal(big.Word(0x07070707), words[7])
}

func TestWrodsToBytes(t *testing.T) {
	require := require.New(t)

	if sizeOfUInt != 4 && sizeOfUInt != 8 {
		t.FailNow()
	}

	raw := make([]big.Word, 8)
	for i := 0; i < 8; i++ {
		raw[i] = big.Word(i)
	}
	bytes := wordsToBytes(raw)
	require.Len(bytes, 8*sizeOfUInt)
	for i := 0; i < 8; i++ {
		bytes[sizeOfUInt*i] = byte(i)
	}
}
