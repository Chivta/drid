package epp

import (
	"encoding/xml"
	"log"
)

const loginExample = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
 <epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <login>
      <clID>test-49</clID>
      <pw>zkXvaw9jHpbjQ4d</pw>
      <options>
        <version>1.0</version>
        <lang>en</lang>
      </options>
      <svcs>
        <objURI>urn:ietf:params:xml:ns:epp-1.0</objURI>
        <objURI>http://hostmaster.ua/epp/contact-1.1</objURI>
        <objURI>http://hostmaster.ua/epp/domain-1.1</objURI>
        <objURI>http://hostmaster.ua/epp/host-1.1</objURI>
        <svcExtension>
          <extURI>http://hostmaster.ua/epp/rgp-1.1</extURI>
          <extURI>http://hostmaster.ua/epp/uaepp-1.1</extURI>
          <extURI>http://hostmaster.ua/epp/balance-1.0</extURI>
          <extURI>http://hostmaster.ua/epp/secDNS-1.1</extURI>
        </svcExtension>
      </svcs>
    </login>
  </command>
 </epp>`

type Request struct {
	XMLName xml.Name `xml:"epp"`
	Command *Command `xml:"command"`
	Greeting *Greeting `xml:"greeting"`
}

type Command struct {
	Login  *LoginCommand  `xml:"login"`
	Logout *LogoutCommand `xml:"logout"`
	Check  *CheckCommand  `xml:"check"`
	Info   *InfoCommand   `xml:"info"`
	Create *CreateCommand `xml:"create"`
	Delete *DeleteCommand `xml:"delete"`
	Renew  *RenewCommand  `xml:"renew"`
	Transfer *TransferCommand `xml:"transfer"`
	Update *UpdateCommand `xml:"update"`
}

type Greeting struct {
	// Define greeting fields if needed
}

type LoginCommand struct {
	ClientID string `xml:"clID"`
	Password string `xml:"pw"`
	Options struct {
		Version string `xml:"version"`
		Lang    string `xml:"lang"`
	} `xml:"options"`
	Services struct {
		ObjURIs []string `xml:"objURI"`
		SvcExtension struct {
			ExtURIs []string `xml:"extURI"`
		} `xml:"svcExtension"`
	} `xml:"svcs"`
}

type LogoutCommand struct{}
type CheckCommand struct{}
type InfoCommand struct{}
type CreateCommand struct{}
type DeleteCommand struct{}
type RenewCommand struct{}
type TransferCommand struct{}
type UpdateCommand struct{}

func HandleRequest(client *Client, request string) (string, error) {
	requestObj := &Request{}
	err := xml.Unmarshal([]byte(request), requestObj)
	if err != nil {
		return "", err
	}

	if requestObj.Command != nil {
		switch {
		case requestObj.Command.Login != nil:
			log.Printf("Login command from client: %s", requestObj.Command.Login.ClientID)
		case requestObj.Command.Logout != nil:
			log.Println("Logout command")
		case requestObj.Command.Check != nil:
			log.Println("Check command")
		case requestObj.Command.Info != nil:
			log.Println("Info command")
		case requestObj.Command.Create != nil:
			log.Println("Create command")
		case requestObj.Command.Delete != nil:
			log.Println("Delete command")
		case requestObj.Command.Renew != nil:
			log.Println("Renew command")
		case requestObj.Command.Transfer != nil:
			log.Println("Transfer command")
		case requestObj.Command.Update != nil:
			log.Println("Update command")
		}
	} else if requestObj.Greeting != nil {
		log.Println("Greeting received")
	}

	return "", nil
}