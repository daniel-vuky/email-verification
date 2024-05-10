package emailverification

// makeResult is the response from the email verification service
func makeResult(email string, valid bool, errorMessage string) *VerifyResponse {
	return &VerifyResponse{
		Email:        email,
		Valid:        valid,
		ErrorMessage: errorMessage,
	}
}
