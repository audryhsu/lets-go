package main

// create a custom string type so ensure that context keys are unique (avoid collison)
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
