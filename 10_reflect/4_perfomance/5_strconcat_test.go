package main

import (
	"bytes"
	"testing"
)

func BenchmarkStrConcatDummy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str := ""
		for j := 0; j < 1000; j++ {
			str += "123"
		}
	}
}

func BenchmarkStrConcatBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buff := bytes.NewBuffer(make([]byte, 3500))
		for j := 0; j < 1000; j++ {
			buff.WriteString("123")
		}
	}
}
