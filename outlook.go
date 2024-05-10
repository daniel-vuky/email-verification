package emailverification

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	SignUpUrl = "https://signup.live.com/signup.aspx?uaid=16&mkt=en-US&lic=1&wreply=https%3a%2f%2faccount.microsoft.com%2faccount%2fCreateAccount%3fru%3dhttps%253a%252f%252faccount.microsoft.com%252faccount%252fCreateAccount&id=64855&cbcxt=mai&bk=1399041236&uiflavor=web&uaid=16&mkt=en-US&lc=1033&lic=1"
)

type Outlook struct {
	client *http.Client
}

func NewOutlookVerifier() MailProvider {
	return &Outlook{
		client: &http.Client{},
	}
}

func (o *Outlook) IsSupportedDomain(domain string) bool {
	return domain == "outlook.com"
}

func (o *Outlook) IsValid(email string) error {
	values := url.Values{}
	values.Set("signupEmail", email)

	resp, err := http.PostForm(fmt.Sprintf("%s&%s", SignUpUrl, values.Encode()), values)
	if err != nil {
		return fmt.Errorf("failed to validate email: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if !strings.Contains(string(body), "The email address you entered is already in use.") {
		return fmt.Errorf("email is not existed")
	}

	return nil
}
