package emailverification

import (
	"errors"
	"golang.org/x/net/idna"
	"net"
	"net/smtp"
	"sync"
	"time"
)

const (
	smtpPort    = ":25"
	smtpTimeout = 30 * time.Second
)

// IsValidSMTP checks if the email is valid
func (v *Verifier) IsValidSMTP(mail, domain string) error {
	smtpClient, _, err := newSMTPClient(domain)
	if err != nil {
		return err
	}
	defer smtpClient.Close()

	if err = smtpClient.Hello(domain); err != nil {
		return err
	}

	if err = smtpClient.Mail(mail); err != nil {
		return err
	}

	if err = smtpClient.Rcpt(mail); err != nil {
		return err
	}

	return nil
}

// newSMTPClient generates a new available SMTP client
func newSMTPClient(domain string) (*smtp.Client, *net.MX, error) {
	domain = domainToASCII(domain)
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, nil, err
	}

	if len(mxRecords) == 0 {
		return nil, nil, errors.New("no MX records found")
	}
	// Create a channel for receiving response from dialing SMTP servers
	ch := make(chan interface{}, 1)
	selectedMXCh := make(chan *net.MX, 1)

	// Done indicates if we're still waiting on dial responses
	var done bool

	// mutex for data race
	var mutex sync.Mutex

	// Attempt to connect to all SMTP servers concurrently
	for i, r := range mxRecords {
		addr := r.Host + smtpPort
		index := i
		go func() {
			c, err := dialSMTP(addr)
			if err != nil {
				if !done {
					ch <- err
				}
				return
			}

			// Place the client on the channel or close it
			mutex.Lock()
			switch {
			case !done:
				done = true
				ch <- c
				selectedMXCh <- mxRecords[index]
			default:
				c.Close()
			}
			mutex.Unlock()
		}()
	}

	// Collect errors or return a client
	var errs []error
	for {
		res := <-ch
		switch r := res.(type) {
		case *smtp.Client:
			return r, <-selectedMXCh, nil
		case error:
			errs = append(errs, r)
			if len(errs) == len(mxRecords) {
				return nil, nil, errs[0]
			}
		default:
			return nil, nil, errors.New("unexpected response dialing SMTP server")
		}
	}

}

func domainToASCII(domain string) string {
	asciiDomain, err := idna.ToASCII(domain)
	if err != nil {
		return domain
	}
	return asciiDomain
}

// dialSMTP is a timeout wrapper for smtp.Dial. It attempts to dial an
// SMTP server (socks5 proxy supported) and fails with a timeout if timeout is reached while
// attempting to establish a new connection
func dialSMTP(addr string) (*smtp.Client, error) {
	// Channel holding the new smtp.Client or error
	ch := make(chan interface{}, 1)

	// Dial the new smtp connection
	go func() {
		var conn net.Conn
		var err error
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			ch <- err
			return
		}

		host, _, _ := net.SplitHostPort(addr)
		client, err := smtp.NewClient(conn, host)
		if err != nil {
			ch <- err
			return
		}
		ch <- client
	}()

	// Retrieve the smtp client from our client channel or timeout
	select {
	case res := <-ch:
		switch r := res.(type) {
		case *smtp.Client:
			return r, nil
		case error:
			return nil, r
		default:
			return nil, errors.New("unexpected response dialing SMTP server")
		}
	case <-time.After(smtpTimeout):
		return nil, errors.New("timeout connecting to mail-exchanger")
	}
}
