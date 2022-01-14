package drizzle

//SecretRetriever handles secrets.
type SecretRetriever interface {
	// RetrieveSecret retrieves a secret for given ID.
	RetrieveSecret(secretID string) (string, error)
}
