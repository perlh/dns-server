package main

import (
	"dns-server/dns"
	"fmt"
)

func main() {
	err := dns.Service()
	if err != nil {
		fmt.Println("dns start fail,error ", err)
		panic(err)
	}
}
