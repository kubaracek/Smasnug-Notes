package cloak

type Cloak interface {
	CloakExecution(func() error) error
}
