// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

// This program generates tables.go:
//	go run maketables.go | gofmt > tables.go

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	fmt.Printf("// generated by go run maketables.go; DO NOT EDIT\n\n")
	fmt.Printf("// Package simplifiedchinese provides Simplified Chinese encodings such as GBK.\n")
	fmt.Printf("package simplifiedchinese\n\n")

	res, err := http.Get("http://encoding.spec.whatwg.org/index-gbk.txt")
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer res.Body.Close()

	mapping := [65536]uint16{}
	reverse := [65536]uint16{}

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" || s[0] == '#' {
			continue
		}
		x, y := uint16(0), uint16(0)
		if _, err := fmt.Sscanf(s, "%d 0x%x", &x, &y); err != nil {
			log.Fatalf("could not parse %q", s)
		}
		if x < 0 || 126*190 <= x {
			log.Fatalf("GBK code %d is out of range", x)
		}
		mapping[x] = y
		if reverse[y] == 0 {
			c0, c1 := x/190, x%190
			if c1 >= 0x3f {
				c1++
			}
			reverse[y] = (0x81+c0)<<8 | (0x40 + c1)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner error: %v", err)
	}

	fmt.Printf("// gbkDecode is the decoding table from GBK code to Unicode.\n")
	fmt.Printf("// It is defined at http://encoding.spec.whatwg.org/index-gbk.txt\n")
	fmt.Printf("var gbkDecode = [...]uint16{\n")
	for i, v := range mapping {
		if v != 0 {
			fmt.Printf("\t%d: 0x%04X,\n", i, v)
		}
	}
	fmt.Printf("}\n\n")

	fmt.Printf("// gbkEncode is the encoding table from Unicode to GBK code.\n")
	fmt.Printf("var gbkEncode = [65536]uint16{\n")
	for i, v := range reverse {
		if v != 0 {
			fmt.Printf("\t%d: 0x%04X,\n", i, v)
		}
	}
	fmt.Printf("}\n\n")
}