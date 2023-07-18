package interfaces

type APCService interface {
	Name() string
	SetName(newValue string)
	IpAddress() string
	SetIpAddress(newValue string)
	Connect() error
	Begin() error
	AddEvent(newValue string)
	AddStatus(newValue string)
	Events() []string
	Status() []string
	PeriodicUpdateStart()
	PeriodicUpdateStop()
	Request(command string, r chan string) error
	SendCommand(command string) error
	ReceiveMessage() (string, error)
	End()
}
