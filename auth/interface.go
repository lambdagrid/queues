package auth

type AuthProvider interface {
	Check(key, secret string) (bool, error)
	CreateAccount(accountName string) (key, secret string, err error)
}
