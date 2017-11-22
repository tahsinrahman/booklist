package main

import (
	"fmt"
	"strconv"
)

type IPAddr [4]byte

func (ip IPAddr) String() string {
	var ret string

	for id, x := range ip {
		if id > 0 {
			ret += "."
		}

		val := strconv.Itoa(int(x))
		ret += val
	}

	return ret
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}

	for name, ip := range hosts {
		fmt.Println(name, ip)
	}

}
