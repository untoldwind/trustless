package secrets

import "github.com/untoldwind/trustless/api"

func (s *Secrets) List() (*api.SecretList, error) {
	if s.IsLocked() {
		return nil, SecretsLockedError
	}
	if err := s.buildIndex(); err != nil {
		return nil, err
	}

	return s.index.list(), nil
}
