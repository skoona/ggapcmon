package interfaces

type Service interface {
	HostMessageChannel(hostName string) chan []string
	Shutdown()
}
