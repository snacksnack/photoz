package email

import (
	"fmt"

	"gopkg.in/mailgun/mailgun-go.v1"
)

const welcomeSubject = "Welcome to Photoz!"

const welcomeText = `Hi there!

Welcome to Photoz! Please contact us with application feedback!

--Reid
`

const welcomeHTML = `Hi there!<br/>
<br/>
Welcome to <a href="https://photoz.reidc.io">Photoz</a>! We hope you enjoy
our application.<br/>
<br/>
--Reid
`

//functional options...
type ClientConfig func(*Client)

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		//set default from address
		from: "support@photoz.reidc.io",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := c.mg.Send(message) //don't care about response or messageid
	return err
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
