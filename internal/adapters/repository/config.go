package repository

import (
	"context"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/core/domain"
	"github.com/skoona/ggapcmon/internal/core/ports"
	"net"
	"time"
)

const (
	HostLocal     = "127.0.0.1:3551"
	HostLocalName = "Local"
	HostsPrefs    = "ApcHost"
)

type config struct {
	hosts map[string]*domain.ApcHost
	prefs fyne.Preferences
}

var _ ports.Configuration = (*config)(nil)
var _ ports.Provider = (*config)(nil)

func NewConfig() (ports.Configuration, error) {
	var err error
	var hosts map[string]*domain.ApcHost
	prefs := fyne.CurrentApp().Preferences()
	defaultHost := domain.NewApcHost(HostLocalName, HostLocal, 10, 5, true, true)
	defaultHosts := map[string]*domain.ApcHost{
		// graph-30 = 15 hours @ 15 network-sec
		defaultHost.Id: defaultHost,
	}

	commons.DebugLog("Default IP: ", commons.DefaultIp())

	hostString := prefs.String(HostsPrefs)
	if hostString != "" && len(hostString) > 16 {
		commons.DebugLog("NewConfig() load ApcHost preferences succeeded ")
		err = json.Unmarshal([]byte(hostString), &hosts)
		if err != nil {
			commons.DebugLog("NewConfig() Unmarshal ApcHhost failed: ", err.Error())
		}
	}
	if len(hosts) == 0 {
		commons.DebugLog("NewConfig() load preferences ApcHost failed using defaults ")
		save, _ := json.Marshal(defaultHosts)
		prefs.SetString(HostsPrefs, string(save))
		hosts = defaultHosts
	}

	// restore binding
	for _, h := range hosts {
		h.Bloadpct = binding.NewString()
		h.Bbcharge = binding.NewString()
		h.Blinev = binding.NewString()
		h.Bcumonbatt = binding.NewString()
		h.Bxoffbatt = binding.NewString()
		h.Bxonbatt = binding.NewString()
		h.Bnumxfers = binding.NewString()
		h.Bmaster = binding.NewString()
		h.Bcable = binding.NewString()
	}

	cfg := &config{
		hosts: hosts,
		prefs: prefs,
	}

	for _, h := range cfg.hosts {
		_ = cfg.VerifyHostConnection(h)
	}

	return cfg, err
}
func (c *config) ResetConfig() {
	c.prefs.SetString(HostsPrefs, "")
}
func (c *config) HostById(id string) *domain.ApcHost {
	return c.hosts[id]
}
func (c *config) Hosts() []*domain.ApcHost {
	var r []*domain.ApcHost
	for _, v := range c.hosts {
		r = append(r, v)
	}
	return r
}
func (c *config) Apply(h *domain.ApcHost) ports.Configuration {
	c.hosts[h.Id] = h
	_ = c.VerifyHostConnection(h)
	return c
}
func (c *config) AddHost(host *domain.ApcHost) {
	c.Apply(host).Save()
	commons.DebugLog("Config::AddHost() saved: .", host)
}
func (c *config) Save() {
	save, err := json.Marshal(c.hosts)
	if err != nil {
		commons.DebugLog("Configuration::Save() marshal apcHosts failed: ", err.Error())
	} else {
		c.prefs.SetString(HostsPrefs, string(save))
	}
}
func (c *config) Remove(id string) {
	if id == "" {
		return
	}
	delete(c.hosts, id)
	c.Save()
}
func (c *config) HostKeys() []string {
	return commons.Keys(c.hosts)
}

// Shutdown compliance with Provider Interface
func (c *config) Close() {
	commons.DebugLog("Config::Close() called.")
}

// VerifyHostConnection compliance with Provider Interface
func (c *config) VerifyHostConnection(h *domain.ApcHost) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", h.IpAddress)
	if err != nil {
		commons.DebugLog("connect() dial Error: ", err.Error(), ", Ip: ", h.IpAddress, ", context: ", ctx.Err())
		h.State = commons.HostStatusUnknown
		if ctx.Err() != nil {
			return ctx.Err()
		} else {
			return err
		}
	}
	time.Sleep(50 * time.Millisecond)
	h.State = commons.HostStatusOnline
	_ = conn.Close()
	return nil
}
