package gen_test

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

func TestSliceOf(t *testing.T) {
	genParams := gopter.DefaultGenParameters()
	genParams.Size = 50
	elementGen := gen.Const("element")
	sliceGen := gen.SliceOf(elementGen)

	for i := 0; i < 100; i++ {
		sample, ok := sliceGen(genParams).Retrieve()

		if !ok {
			t.Error("Sample was not ok")
		}
		strings, ok := sample.([]string)
		if !ok {
			t.Errorf("Sample not slice of string: %#v", sample)
		} else {
			if len(strings) >= 50 {
				t.Errorf("Sample has invalid length: %#v", len(strings))
			}
			for _, str := range strings {
				if str != "element" {
					t.Errorf("Sample contains invalid value: %#v", sample)
				}
			}
		}
	}

	genParams.Size = 0

	for i := 0; i < 100; i++ {
		sample, ok := sliceGen(genParams).Retrieve()

		if !ok {
			t.Error("Sample was not ok")
		}
		strings, ok := sample.([]string)
		if !ok {
			t.Errorf("Sample not slice of string: %#v", sample)
		} else {
			if len(strings) != 0 {
				t.Errorf("Sample has invalid length: %#v", len(strings))
			}
		}
	}
}

func TestSliceOfN(t *testing.T) {
	elementGen := gen.Const("element")
	sliceGen := gen.SliceOfN(10, elementGen)

	for i := 0; i < 100; i++ {
		sample, ok := sliceGen.Sample()

		if !ok {
			t.Error("Sample was not ok")
		}
		strings, ok := sample.([]string)
		if !ok {
			t.Errorf("Sample not slice of string: %#v", sample)
		} else {
			if len(strings) != 10 {
				t.Errorf("Sample has invalid length: %#v", len(strings))
			}
			for _, str := range strings {
				if str != "element" {
					t.Errorf("Sample contains invalid value: %#v", sample)
				}
			}
		}
	}
}

func TestSliceOfNSieve(t *testing.T) {
	var called int
	elementSieve := func(v interface{}) bool {
		called++
		return v == "element"
	}
	elementGen := gen.Const("element").SuchThat(elementSieve)
	sliceGen := gen.SliceOfN(10, elementGen)
	result := sliceGen(gopter.DefaultGenParameters())
	value, ok := result.Retrieve()
	if !ok || value == nil {
		t.Errorf("Invalid value: %#v", value)
	}
	strs, ok := value.([]string)
	if !ok || len(strs) != 10 {
		t.Errorf("Invalid value: %#v", value)
	}
	if called != 20 {
		t.Errorf("Invalid called: %d", called)
	}
	if result.Sieve(strs[0:9]) {
		t.Error("Sieve must not allow array len < 10")
	}
	strs[0] = "bla"
	if result.Sieve(strs) {
		t.Error("Sieve must not allow array with invalid element")
	}
}
