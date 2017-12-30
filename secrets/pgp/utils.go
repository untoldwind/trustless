package pgp

import (
	"math/big"
	"reflect"
	"unsafe"

	"github.com/awnumar/memguard"
	"github.com/pkg/errors"
)

type writeTo struct {
	buffer   *memguard.LockedBuffer
	position int
}

func (w *writeTo) Write(p []byte) (n int, err error) {
	if len(p)+w.position > w.buffer.Size() {
		return 0, errors.New("Overflow")
	}
	if err := w.buffer.MoveAt(p, w.position); err != nil {
		return 0, errors.Wrap(err, "Memcopy")
	}
	w.position += len(p)
	return len(p), nil
}

func (w *writeTo) Result() []byte {
	w.buffer.MakeImmutable()
	return w.buffer.Buffer()[0:w.position]
}

var sizeOfUInt int

func bytesToWords(raw []byte) []big.Word {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len /= sizeOfUInt
	header.Cap /= sizeOfUInt
	return *(*[]big.Word)(unsafe.Pointer(&header))
}

func wordsToBytes(raw []big.Word) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= sizeOfUInt
	header.Cap *= sizeOfUInt
	return *(*[]byte)(unsafe.Pointer(&header))
}

func init() {
	a := uint(0)

	sizeOfUInt = int(unsafe.Sizeof(a))
}
