package main

import (
	"github.com/reiver/go-telnet"
)

func main() {
	handler := AsciiHandler{}
	err := telnet.ListenAndServe(":23", handler)
	if err != nil {
		panic(err)
	}
}
