package filesystem

type SecretManager struct{}

func NewSecretManager() *SecretManager {
	return &SecretManager{}
}

func (p *SecretManager) GetSecret(env string, id string) (string, error) {
	return id, nil
}
