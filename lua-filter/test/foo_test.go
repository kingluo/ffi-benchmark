package test

import "testing"

func BenchmarkRAW(t *testing.B) {
	for i := 0; i < 1000000; i++ {
		testRAW(t)
	}
}

func BenchmarkJSON(t *testing.B) {
	for i := 0; i < 1000000; i++ {
		testJSON(t)
	}
}
