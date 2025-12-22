package main

import (
	"drid/internal/epp"
	"drid/internal/whois"
)

func main() {
	go epp.ListenForEPPConnections(":7000")
	go whois.ListenForWhoisConnections(":4300")

	select {}
}
