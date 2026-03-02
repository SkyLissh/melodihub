package spotify

import (
	"encoding/json"
	"time"
)

type tokenExpiry time.Time

func (t *tokenExpiry) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	now := time.Now()
	*t = tokenExpiry(now.Add(time.Duration(value) * time.Second))
	return nil
}

type Token struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresAt   tokenExpiry `json:"expires_in"`
	Scope       string      `json:"scope"`
}

func (t *Token) IsExpired() bool {
	now := time.Now()
	skew := 30 * time.Second
	return now.Add(skew).After(time.Time(t.ExpiresAt))
}
