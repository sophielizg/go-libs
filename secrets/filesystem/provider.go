package filesystem

type SecretProvider struct{}

func NewSecretProvider() *SecretProvider {
	return &SecretProvider{}
}

func (p *SecretProvider) GetSecret(env string, id string) (string, error) {
	return id, nil
}
