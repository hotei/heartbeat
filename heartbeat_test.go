// heartbeat_test.go (c) 2010-2014 David Rook - all rights reserved

package main

import (
	"fmt"
	"net"
	"testing"
	"time"
	// local
	"github.com/hotei/mdr"
)

func Test_0001(t *testing.T) {
	// fmt.Printf("go test -bench=\".*\" to run all benchmarks\n")
	// fmt.Printf("to run single test E use go test -run=\"Test_E\"\n")

	var sec int64 = 86400 * 600 // 600 days

	rcs := mdr.HumanTime(time.Duration(sec * 1000 * 1000 * 1000))
	fmt.Printf("Human readable time for (600 days) %d sec is : %s\n", sec, rcs)

	teststr := "10.1.2.115"
	x := net.ParseIP(teststr)
	n := mdr.Uint32FromIP(x)
	fmt.Printf("Teststr: %s=>  IP(%v) mdr.Int32FromIP(%08x)\n", teststr, x, n)

	nx := mdr.IPFromUint32(n)
	fmt.Printf("nx %v\n", nx)
}
