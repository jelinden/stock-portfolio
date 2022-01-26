package email

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jelinden/stock-portfolio/app/config"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

var mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

var domain = os.Getenv("MAILGUN_DOMAIN")
var privateAPIKey = os.Getenv("MAILGUN_API_KEY")

func SendVerificationEmail(emailTo, hash, fromEmail string) {
	message := "<div>Please verify your account with following link:<br/>" +
		"<a href=\"" + config.Config.VerifyURL + hash + "\">Verify</a>.<br/><br/>" +
		"If you received this message without registering to portfolio.jelinden.fi," +
		" you can delete the message.</div>"
	subject := "Please verify your new account"

	sendEmail(emailTo, fromEmail, subject, message)
}

func sendEmail(emailTo, fromEmail, subject, message string) {
	fmt.Println("-", privateAPIKey)
	mg := mailgun.NewMailgun(domain, privateAPIKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	msg := mg.NewMessage(fromEmail, subject, message, emailTo)
	msg.SetHtml(message)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
