package main

import (
	"drid/internal/config"
	"drid/internal/epp"
	"drid/internal/whois"
	"drid/pkg/db"
	"log"
)

func main() {
	cfg,err:= config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Could not load config: %v\n",err)
	}

	db,err := db.NewDBConnection(cfg.DB.Host,cfg.DB.User,cfg.DB.Password,cfg.DB.Name,cfg.DB.Port)
	if err != nil {
		log.Fatalf("Could not create db connection: %v\n",err)
	}
	defer db.Close()

	whoisServer := whois.NewWhoisServer(db)
	eppServer := epp.NewEppServer(db)
	go whoisServer.ListenForWhoisConnections(cfg.Server.WhoisAddress)
	go eppServer.ListenForEPPConnections(cfg.Server.EPPAddress)

	select {}
}
