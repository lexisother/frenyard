// Inspired by go-bindata, but that seemed too heavyweight.

package main

import (
	"flag"
	"fmt"
	"os"
	"io"
	"encoding/base64"
)

func main() {
	flag.Parse()
	target, err := os.OpenFile("data.go", os.O_WRONLY | os.O_TRUNC | os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(target, "package design")
	fmt.Fprintln(target, "// Generated code; update using: go generate ./... in CCUpdaterUI main dir.")
	for x := 0; x < flag.NArg(); x++ {
		name := flag.Arg(x)
		file, err := os.Open("bindata/" + name + ".png")
		if err != nil {
			panic(err)
		}
		stat, err := file.Stat()
		if err != nil {
			panic(err)
		}
		data := make([]byte, stat.Size())
		amount, err := io.ReadFull(file, data)
		if err != nil || amount != len(data) {
			panic(err)
		}
		file.Close()
		fmt.Fprintln(target, "var " + name + "B64 = \"" + base64.StdEncoding.EncodeToString(data) + "\"")
	}
	target.Close()
}
