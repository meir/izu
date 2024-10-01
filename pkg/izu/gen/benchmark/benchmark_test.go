//go:generate go run ./gen.go
package main

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func get_key() string {
	size := len(keys_array)
	return keys_array[rand.Intn(size)]
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validate_map(get_key())
	}
}

func BenchmarkSwitch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validate_switch(get_key())
	}
}

func BenchmarkArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validate_array(get_key())
	}
}
