package emailverification

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const gmailAPI = "https://mail.google.com/mail/gxlu?email=%s"

type Gmail struct {
	client *http.Client
}

func NewGmailVerifier() MailProvider {
	return &Gmail{
		client: &http.Client{},
	}
}

func (g *Gmail) IsSupportedDomain(domain string) bool {
	return domain == "gmail.com"
}

func (g *Gmail) IsValid(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(gmailAPI, email), nil)
	if err != nil {
		return err
	}
	response, err := g.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to validate email")
	}
	if len(response.Cookies()) == 0 {
		return fmt.Errorf("email is not valid")
	}

	return nil
}
