package shorten

import "testing"

func BenchmarkURLShorten(b *testing.B) {
	for i := 0; i < b.N; i++ {
		URLShorten()
	}
}
