package config

import "time"

type ContextKey string

// Context key for storing a *domain.User in request context.
const UserContextKey = ContextKey("user")

// Default time.Duration for request timeout
const DefaultTimeout = 5 * time.Second
