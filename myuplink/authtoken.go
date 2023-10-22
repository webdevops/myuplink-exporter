package myuplink

import (
	"time"
)

type (
	AuthToken struct {
		AccessToken string `json:"access_token"`
		expiry      time.Time
		ExpiresIn   float64 `json:"expires_in"`
		TokenType   string  `json:"token_type"`
		Scope       string  `json:"scope"`
	}
)

func (t *AuthToken) Process() *AuthToken {
	// reduce expiry time to ensure refresh right before
	expirySeconds := t.ExpiresIn - (t.ExpiresIn * 0.1)
	t.expiry = time.Now().Add(time.Duration(expirySeconds) * time.Second)

	return t
}

func (t *AuthToken) IsExpired() bool {
	return time.Now().After(t.expiry)
}
