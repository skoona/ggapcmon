package entities

import "time"

type ApcHost struct {
	IpAddress        string
	Name             string
	SecondsPerSample time.Duration
}

func NewApcHost(name, ip string, secondsPerSample time.Duration) ApcHost {
	return ApcHost{
		IpAddress:        ip,
		Name:             name,
		SecondsPerSample: secondsPerSample,
	}
}
