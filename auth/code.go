package auth

import (
	"github.com/unickorn/golem-poll-manager/internal/roundtrip"
	"math/rand"
	"net/http"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var src rand.Source

func init() {
	// nanosecond source - was not working with tinygo for some reason :(
	//src = rand.NewSource(int64(out.WasiClocksWallClockNow().Seconds*uint64(math.Pow10(9)) + uint64(out.WasiClocksWallClockNow().Nanoseconds)))

	// use hand-picked seed for now
	src = rand.NewSource(247889329438375)
	http.DefaultClient.Transport = roundtrip.WasiHttpTransport{}
}

// RandStringBytesRmndr returns a random string of length n.
// credit: https://stackoverflow.com/a/31832326/12308931
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[src.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
