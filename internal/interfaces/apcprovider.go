package interfaces

type ApcProvider interface {
	Name() string
	IpAddress() string
	Shutdown()
}
