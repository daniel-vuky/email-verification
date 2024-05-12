package emailverification

import "fmt"

type MailProvider interface {
	IsSupportedDomain(domain string) bool
	IsValid(email string) error
}

func (v *Verifier) checkWithMailProvider(email, provider string) error {
	if v.ListProviderVerifier[provider] == nil || !v.ListProviderVerifier[provider].IsSupportedDomain(provider) {
		return fmt.Errorf("provider %s is not supported", provider)
	}

	return v.ListProviderVerifier[provider].IsValid(email)
}
