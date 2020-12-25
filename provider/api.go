package provider

type API interface {
	GetSecret(key string) (string, error)
	IsSecret(key string) bool
}
