package db

import (
	"database/sql"
	"fmt"
)

func (db *DB) Whois(domain string) (string, error) {
	var info string
	xml := fmt.Sprintf(`<whois type="normal" obj="%s" attr=""><name>%s</name></whois>`, "domain", domain)
	query := fmt.Sprintf(`SELECT uanic_data.whois_exec(E'%s'::text);`,xml)
	
	err := db.conn.QueryRow(query).Scan(&info)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil 
		}
		return "", err
	}
	return info, nil
}
