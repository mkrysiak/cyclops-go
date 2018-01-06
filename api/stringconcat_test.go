package api

import (
	"bytes"
	"testing"
)

func BenchmarkStringConcat(b *testing.B) {
	var str string
	for n := 0; n < b.N; n++ {
		str += "httprequest"
	}
	b.StopTimer()
}

func BenchmarkBufferConcat(b *testing.B) {
	var buffer bytes.Buffer
	for n := 0; n < b.N; n++ {
		buffer.WriteString("httprequest")
	}
	b.StopTimer()
}
