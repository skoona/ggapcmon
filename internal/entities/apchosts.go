package entities

import (
	"github.com/skoona/ggapcmon/internal/commons"
	"time"
)

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
		State:                commons.HostStatusUnknown,
	}
}
func (a *ApcHost) IsEmpty() bool {
	return a.Name != ""
}
