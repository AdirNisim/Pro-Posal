package models

import (
	"time"

	"github.com/google/uuid"
)

// POST /auth/token -> To obtain token by email and password -> Bearer Token + Refresh Token
// POST /auth/refresh -> To refresh a token and extend its expiration time, using a refresh token -> Bearer Token + Refresh Token

// Option 1 - Don't store it anywhere, just encode it and pass back and forth with the caller
// - Pros: No need to store it, lower latencies and storage costs.
// - Cons: You "lose" some of the control over it, as you cannot really revoke it on your side.

// Option 2 - Store it somewhere, and pass as part of the token only the sessionID
// - Pros: You can revoke it, and you can have more control over it, also the token is smaller in size.
// - Const: You need to manage the persistent layer.

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AuthToken struct {
	BearerToken  string    `json:"bearer_token"`
	// RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
