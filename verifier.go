package emailverification

import (
	"strings"
)

type Verifier struct {
	ListProviderVerifier map[string]MailProvider
	CheckMailProvider    bool
}

type VerifyResponse struct {
	Email        string
	Valid        bool
	ErrorMessage string
}

func NewVerifier(checkMailProvider bool) *Verifier {
	return &Verifier{
		CheckMailProvider: checkMailProvider,
		ListProviderVerifier: map[string]MailProvider{
			"gmail":   NewGmailVerifier(),
			"outlook": NewOutlookVerifier(),
		},
	}
}

func (v *Verifier) Verify(email string) *VerifyResponse {
	var err error
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return makeResult(email, false, "not valid RFC format")
	}

	if !v.IsValidEmail(email) {
		return makeResult(email, false, "invalid email format")
	}

	err = v.IsValidSMTP(email, parts[1])
	if err != nil {
		return makeResult(email, false, err.Error())
	}

	err = v.checkWithMailProvider(email, parts[1])
	if err != nil {
		return makeResult(email, false, err.Error())
	}

	return makeResult(email, true, "")
}
