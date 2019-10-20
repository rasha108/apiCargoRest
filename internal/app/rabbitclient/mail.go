package rabbitclient

import "net/mail"

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
