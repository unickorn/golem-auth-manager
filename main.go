package main

import (
	"fmt"
	"github.com/unickorn/golem-poll-manager/auth"
	poll "github.com/unickorn/golem-poll-manager/manager_out"
	"github.com/unickorn/golem-poll-manager/verify"
)

func main() {
}

func init() {
	i := GolemPollManager{}
	poll.SetExportsIlhanGolemPollManagerApi(i)

	// Mailjet API Key and Secret
	auth.Initialize("xxxxxxx", "xxxxxxx")
	// JWT HMAC Secret (unfortunately symmetric for now, RSA is not fully supported yet on tinygo)
	verify.Initialize("super-secret-hmac-secret")
}

// GolemPollManager is the implementation of the manager interface of Golem Poll, a poll template for Golem Cloud.
type GolemPollManager struct {
}

// Login sends an authentication email to the user.
func (g GolemPollManager) Login(email string) poll.Result[string, string] {
	err := auth.SendAuthenticationEmail(email)
	if err != nil {
		return poll.Result[string, string]{
			Kind: poll.Err,
			Err:  fmt.Sprintf("Failed to send authentication email: %s", err.Error()),
		}
	}
	return poll.Result[string, string]{
		Kind: poll.Ok,
		Val:  fmt.Sprintf("Authentication email sent to %s!", email),
	}
}

// Verify verifies the user's secret code and returns the access and refresh tokens.
func (g GolemPollManager) Verify(code string) poll.Result[string, string] {
	if mail, ok := auth.VerifyAuthentication(code); ok {
		tokens, err := verify.CreateTokens(mail)
		if err != nil {
			return poll.Result[string, string]{
				Kind: poll.Err,
				Err:  fmt.Sprintf("Failed to create tokens: %s", err.Error()),
			}
		}
		return poll.Result[string, string]{
			Kind: poll.Ok,
			Val:  tokens,
		}
	}
	return poll.Result[string, string]{
		Kind: poll.Err,
		Err:  "Invalid code, authentication failed!",
	}
}

// Refresh refreshes the user's access token.
func (g GolemPollManager) Refresh(email string, refreshToken string) poll.Result[string, string] {
	tokens, err := verify.RefreshToken(email, refreshToken)
	if err != nil {
		return poll.Result[string, string]{
			Kind: poll.Err,
			Err:  fmt.Sprintf("Failed to refresh tokens: %s", err.Error()),
		}
	}
	return poll.Result[string, string]{
		Kind: poll.Ok,
		Val:  tokens,
	}

}

// Logout invalidates the user's refresh token.
func (g GolemPollManager) Logout(accessToken string) poll.Result[struct{}, string] {
	err := verify.InvalidateToken(accessToken)
	if err != nil {
		return poll.Result[struct{}, string]{
			Kind: poll.Err,
			Err:  fmt.Sprintf("Failed to invalidate token: %s", err.Error()),
		}
	}
	return poll.Result[struct{}, string]{
		Kind: poll.Ok,
		Val:  struct{}{},
	}
}

// Validate validates the user's access token.
func (g GolemPollManager) Validate(accessToken string) bool {
	return verify.CheckToken(accessToken) == nil
}
