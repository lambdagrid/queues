package auth

type AuthProvider interface {
	Check(key, secret string) (bool, error)
}
