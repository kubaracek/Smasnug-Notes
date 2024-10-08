package cloak

type Cloak interface {
	CloakExecuting(func() error) error
}
