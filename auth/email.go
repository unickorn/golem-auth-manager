package auth

import (
	"fmt"
	"time"
)

// codeToMail is a map of secret codes to email addresses.
var codeToMail = make(map[string]string)

// emailToCode is a map of email addresses to secret codes.
var emailToCode = make(map[string]string)

func generateCodeFor(email string) string {
	var code string
	for {
		code = RandStringBytesRmndr(8)
		if _, ok := codeToMail[code]; !ok {
			break
		}
		fmt.Println("Duplicate code:", code)
	}
	emailToCode[email] = code
	return code
}

// SendAuthenticationEmail sends an email to the user with their secret code.
func SendAuthenticationEmail(email string) error {
	code := generateCodeFor(email)
	fmt.Printf("Sending code %s to %s\n", code, email)

	err := SendMailjetEmail(email, code)
	if err != nil {
		return err
	}
	codeToMail[code] = email
	emailToCode[email] = code

	// delete the code after 5 minutes
	go func() {
		fmt.Printf("Deleting code in 5 minutes for %s\n", email)
		time.Sleep(5 * time.Minute)
		fmt.Printf("Deleting code now for %s\n", email)
		delete(codeToMail, code)
		delete(emailToCode, email)
	}()
	return nil
}

// VerifyAuthentication verifies the user's secret code and returns the email address if it is valid.
func VerifyAuthentication(code string) (string, bool) {
	if mail, ok := codeToMail[code]; ok {
		delete(codeToMail, code)
		delete(emailToCode, mail)
		return mail, true
	}
	return "", false
}
