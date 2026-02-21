package email

import (
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/ynori7/music/config"
)

type Mailer struct {
	config      config.Config
	emailClient *mailjet.Client
}

func NewMailer(conf config.Config) Mailer {
	return Mailer{
		config:      conf,
		emailClient: mailjet.NewMailjetClient(conf.Email.PublicKey, conf.Email.PrivateKey),
	}
}

func (m Mailer) SendMail(subject string, htmlBody string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.config.Email.From.Address,
				Name:  m.config.Email.From.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: m.config.Email.To.Address,
					Name:  m.config.Email.To.Name,
				},
			},
			Subject:  subject,
			HTMLPart: htmlBody,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.emailClient.SendMailV31(&messages)

	return err
}
