package secrets

type Provider interface {
	GetSecret(env string, id string) (string, error)
}
