package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var Mock bool = true

func Use(args ...any) {}

func MustMarshal(v any) []byte {
	jsonData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonData
}

func MustUnmarshal(b []byte, a any) {
	if err := json.Unmarshal(b, a); err != nil {
		panic(err)
	}
}

func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func MustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// todo: this can do?
func MustGetFromMap(m map[string]any, key string) any {
	ret, found := m[key]
	if !found {
		panic("MustGetFromMap fail")
	}
	return ret
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numRunes = []rune("0123456789")
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numRunes[r.Intn(len(numRunes))]
	}
	return string(b)
}

func MustReadAll(path string) string {
	data, err := os.ReadFile(path) // os.ReadFile handles open, read, and close.
	if err != nil {
		panic(fmt.Sprintf("MustReadAll failed to read file '%s': %v", path, err))
	}
	return string(data)
}
