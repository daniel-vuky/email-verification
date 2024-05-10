**Email Verification**

This is a simple Go package for checking if is valid and can received emails.

**Features**

- Check if an email address is valid and can receive emails
- Supports Outlook email addresses
- Supports Gmail addresses

**Requirements**

- Go 1.13 or higher
- A valid email address with access to SMTP server


**Installation**

To install the package, run the following command:
```
go get -u github.com/daniel-vuky/email-verification
```


Here is an example of how to use the package:

```go
package main

import (
	"fmt"

	emailVerification "github.com/daniel-vuky/email-verification/email-verification"
)

func main() {
	// Use the email-verification package
	// to verify the email address
	email := "daniel.vuky@gmail.com"
	selfVerifier := emailVerification.NewVerifier(true)
	res := selfVerifier.Verify(email)
	fmt.Println("email validation result", res)
	return
}

```

**Contributing**

Contributions are welcome! Please submit a pull request with your changes.

**License**

This project is licensed under the MIT License - see the LICENSE.md file for details.