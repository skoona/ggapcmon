package entities

import "time"

type ApcHost struct {
	IpAddress            string
	Name                 string
	NetworkSamplePeriod  time.Duration
	GraphingSamplePeriod time.Duration
	Enabled              bool
	TrayIcon             bool
	State                string
}

func NewApcHost(name, ip string, networkSamplePeriod, graphingSamplePeriod time.Duration, enable, trayIcon bool) *ApcHost {
	return &ApcHost{
		IpAddress:            ip,
		Name:                 name,
		NetworkSamplePeriod:  networkSamplePeriod,
		GraphingSamplePeriod: graphingSamplePeriod,
		Enabled:              enable,
		TrayIcon:             trayIcon,
		State:                "unknown", // unknown, onbatt, charging, online, unplugged
	}
}
