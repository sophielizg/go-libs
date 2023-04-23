package secrets

import (
	"github.com/sophielizg/go-libs/config"
)

type settings struct {
	Env      string
	Dir      string
	Provider Provider
}

func New(options ...func(*settings)) (map[string]string, error) {
	s := &settings{}

	for _, option := range options {
		option(s)
	}

	secrets := map[string]string{}
	err := config.New(secrets, config.WithEnv(s.Env), config.WithDir(s.Dir))
	if err != nil {
		return nil, err
	}

	for key, id := range secrets {
		secrets[key], err = s.Provider.GetSecret(s.Env, id)
		if err != nil {
			return nil, err
		}
	}

	return secrets, nil
}

func WithEnv(env string) func(*settings) {
	return func(s *settings) {
		s.Env = env
	}
}

func WithDir(dir string) func(*settings) {
	return func(s *settings) {
		s.Dir = dir
	}
}

func WithProvider(provider Provider) func(*settings) {
	return func(s *settings) {
		s.Provider = provider
	}
}
