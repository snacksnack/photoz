package email

import (
	"fmt"
	"net/url"

	"gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to Photoz!"
	resetSubject   = "Instructions for resetting your password"
	resetBaseURL   = "https://photoz.reidc.io/reset"
)

const welcomeText = `Hi there!

Welcome to Photoz! Please contact us with application feedback!

--Photoz Support
`

const welcomeHTML = `Hi there!<br/>
<br/>
Welcome to <a href="https://photoz.reidc.io">Photoz</a>! We hope you enjoy
our application.<br/>
<br/>
Photoz Support
`

const resetTextTmpl = `Hi there!

It appears that you have requested a password reset. If this was you, please follow the
link below to update your password:

%s

If you are asked for a token, please use the following value:

%s

If you did not request a password reset, you can ignore this email and your account will
remain unchanged.

Best,
Photoz Support
`

const resetHTMLTmpl = `Hi there!<br/>
<br/>
It appears that you have requested a password reset. If this was you, please follow the
link below to update your password:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token, please use the following value:<br/>
<br/>
%s<br/>
<br/>
If you did not request a password reset, you can ignore this email and your account will
remain unchanged.<br/>
<br/>
Best,<br/>
Photoz Support<br/>
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

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(c.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)
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
