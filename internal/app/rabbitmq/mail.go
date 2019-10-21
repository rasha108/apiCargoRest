package rabbitmq

import (
	"net/mail"
	"strings"
)

//
type Mail struct {
	From        *mail.Address   `json:"from,omitempty"`
	To          []*mail.Address `json:"to,omitempty"`
	Cc          []*mail.Address `json:"cc,omitempty"`
	Bcc         []*mail.Address `json:"bcc,omitempty"`
	Subject     string          `json:"subject,omitempty"`
	Text        string          `json:"text,omitempty"`
	Attachments []Attachment    `json:"attachments,omitempty"`
}

type Attachment struct {
	Name    string `json:"name,omitempty"`    // File name
	Content string `json:"content,omitempty"` // Base64 encoded []byte
}

func NewSimpleMail(from string, to []string, subject, text string) (*Mail, error) {
	fromAddr, err := mail.ParseAddress(from)
	if err != nil {
		return nil, err
	}

	toAddr, err := mail.ParseAddressList(strings.Join(to, ","))
	if err != nil {
		return nil, err
	}
	m := &Mail{
		From:    fromAddr,
		To:      toAddr,
		Subject: subject,
		Text:    text,
	}

	return m, nil
}
