package config

type ContextKey string

// UserContextKey is the context key for storing a *domain.User in request context.
const UserContextKey = ContextKey("user")
