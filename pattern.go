package emailverification

import "regexp"

const regexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func (v *Verifier) IsValidEmail(email string) bool {
	re := regexp.MustCompile(regexPattern)
	return re.MatchString(email)
}
