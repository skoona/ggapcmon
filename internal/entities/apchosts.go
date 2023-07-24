package entities

import (
	"fyne.io/fyne/v2/data/binding"
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
	Bloadpct             binding.String `json:"-"`
	Bbcharge             binding.String `json:"-"`
	Blinev               binding.String `json:"-"`
	Bcumonbatt           binding.String `json:"-"`
	Bxoffbatt            binding.String `json:"-"`
	Blastxfer            binding.String `json:"-"`
	Bnumxfers            binding.String `json:"-"`
	Bstatus              binding.String `json:"-"`
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
		Bloadpct:             binding.NewString(),
		Bbcharge:             binding.NewString(),
		Blinev:               binding.NewString(),
		Bcumonbatt:           binding.NewString(),
		Bxoffbatt:            binding.NewString(),
		Blastxfer:            binding.NewString(),
		Bnumxfers:            binding.NewString(),
		Bstatus:              binding.NewString(),
	}
}
func (a *ApcHost) IsNil() bool {
	return (a == nil)
}
