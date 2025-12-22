package db

import (
	"database/sql"
)

func (db *DB) Whois(domain string) (string, error) {
	var info string
	err := db.conn.QueryRow(`SELECT uanic_data.whois_exec(E'<whois type="normal" obj="$1" attr=""><name>$2</name></whois>'::text)`,
	"domain", domain).Scan(&info)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil 
		}
		return "", err
	}
	return info, nil
}
