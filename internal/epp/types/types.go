package types

import (
	"encoding/xml"
)

type EPPRequest struct {
	XMLName xml.Name `xml:"epp"`

	Command  *Command  `xml:"command"`
	Greeting *Greeting `xml:"greeting"`
}

type Command struct {
	XMLName xml.Name `xml:"command"`

	Login    *LoginCommand    `xml:"login"`
	Logout   *LogoutCommand   `xml:"logout"`
	Check    *CheckCommand    `xml:"check"`
	Info     *InfoCommand     `xml:"info"`
	Create   *CreateCommand   `xml:"create"`
	Delete   *DeleteCommand   `xml:"delete"`
	Renew    *RenewCommand    `xml:"renew"`
	Transfer *TransferCommand `xml:"transfer"`
	Update   *UpdateCommand   `xml:"update"`
	ClTRID   string           `xml:"clTRID,omitempty"`
}

type Greeting struct {
	XMLName    xml.Name `xml:"greeting"`
	ServerID   string   `xml:"svID"`
	ServerDate string   `xml:"svDate"`
}

type LoginCommand struct {
	ClientID string `xml:"clID"`
	Password string `xml:"pw"`
	Options  struct {
		Version string `xml:"version"`
		Lang    string `xml:"lang"`
	} `xml:"options"`
	Services struct {
		ObjURIs      []string `xml:"objURI"`
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

type EPPResponse struct {
	XMLName xml.Name `xml:"epp"`
	Response *Response `xml:"response"`
	Greeting *Greeting `xml:"greeting"`
}

type Response struct {
	XMLName xml.Name `xml:"response"`

	Result  []Result `xml:"result"`
	TrID    *TrID    `xml:"trID"`
}

type Result struct {
	Code int `xml:"code,attr"`
	Msg  Msg `xml:"msg"`
}

type Msg struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",chardata"`
}

type TrID struct {
	ClTRID string `xml:"clTRID,omitempty"`
	SvTRID string `xml:"svTRID"`
}

type Client struct {
	ClientID      string
	Authenticated bool
}

