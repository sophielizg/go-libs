package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type settings struct {
	Env string
	Dir string
}

func New[T any](cfg T, options ...func(*settings)) error {
	s := &settings{}

	for _, option := range options {
		option(s)
	}

	configBytes, err := os.ReadFile(s.Dir + "/" + s.Env + ".yaml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(configBytes, cfg)
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
