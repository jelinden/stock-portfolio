package email

import (
	"os"
	"testing"

	"github.com/jelinden/stock-portfolio/app/config"
)

func TestEmailSending(t *testing.T) {
	domain = os.Getenv("MAILGUN_DOMAIN")
	privateAPIKey = os.Getenv("MAILGUN_API_KEY")
	config.Config.VerifyURL = "test"
	SendVerificationEmail(os.Getenv("ADMINUSER"), "testing", os.Getenv("FROMEMAIL"))
}
