package persistence

type contextKey string

const (
	dataStoreKey contextKey = "ds"
	txKey        contextKey = "tx"
)
