package secrets

type Manager interface {
	GetSecret(env string, id string) (string, error)
}
