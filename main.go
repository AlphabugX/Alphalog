package main

import (
	"Alphalog/Config"
	"Alphalog/Service"
	"fmt"
)

func main() {
	fmt.Println(Service.Banner)
	Alphalog()
}
func Alphalog()  {
	Config.Initialization()
	Domain := Config.Init.Domain
	// HTTP	Service	Started
	go Service.Httpserver()
	// JNDI Service Started
	go Service.JNDI()
	// DNS	Service Started
	Service.Dnsserver(Domain)
}