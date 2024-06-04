package shared

import "errors"

var (
	ErrBoxNotEligible     = errors.New("box not eligible")
	ErrBoxNotYetAvailable = errors.New("box not yet available")
	ErrNoWebhookUrl       = errors.New("empty webhook url")
	ErrNoAccessToken      = errors.New("no access token")
)
