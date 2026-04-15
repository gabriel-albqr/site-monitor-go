package domain

import "time"

// CheckStatus representa o estado final de uma checagem.
type CheckStatus string

const (
	// CheckStatusOnline indica que o site respondeu dentro do esperado.
	CheckStatusOnline CheckStatus = "online"
	// CheckStatusOffline indica que o site não respondeu ou retornou erro.
	CheckStatusOffline CheckStatus = "offline"
)

// CheckResult representa o resultado de uma checagem de disponibilidade.
type CheckResult struct {
	Site         Site
	HTTPStatus   int
	Status       CheckStatus
	ResponseTime time.Duration
	CheckedAt    time.Time
	ErrorMessage string
}
