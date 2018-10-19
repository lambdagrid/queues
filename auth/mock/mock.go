package mock

import "github.com/lambdagrid/queues/auth"

type MockAuthProvider struct{}

func (m *MockAuthProvider) Check(key, secret string) (bool, error) {
	return true, nil
}

func (m *MockAuthProvider) CreateAccount(accountName string) (key, secret string, err error) {
	return "", "", nil
}

func New() auth.AuthProvider {
	return &MockAuthProvider{}
}
