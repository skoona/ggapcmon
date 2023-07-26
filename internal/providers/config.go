package providers

import (
	"context"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"net"
	"time"
)

const (
	HostLocal     = "127.0.0.1:3551"
	HostLocalName = "Local"
	HostsPrefs    = "ApcHost"
	HubHostsPrefs = "HubHost"
)

type config struct {
	hosts map[string]*entities.ApcHost
	hubs  []*entities.HubHost
	prefs fyne.Preferences
}

var _ interfaces.Configuration = (*config)(nil)
var _ interfaces.Provider = (*config)(nil)

func NewConfig(prefs fyne.Preferences) (interfaces.Configuration, error) {
	var err error
	var hosts map[string]*entities.ApcHost
	var hubHosts []*entities.HubHost

	defaultHosts := map[string]*entities.ApcHost{
		// graph-30 = 15 hours @ 15 network-sec
		HostLocalName: entities.NewApcHost(HostLocalName, HostLocal, 10, 5, true, true),
	}
	defaultHubHosts := []*entities.HubHost{
		entities.NewHubHost("Scotts", "10.100.1.41", "a79c07db-9178-4976-bd10-428aa0d3d159", "10.100.1.183"),
	}

	commons.DebugLog("Default IP: ", commons.DefaultIp())

	hostString := prefs.String(HostsPrefs)
	if hostString != "" && len(hostString) > 16 {
		commons.DebugLog("NewConfig() load ApcHost preferences succeeded ")
		err = json.Unmarshal([]byte(hostString), &hosts)
		if err != nil {
			commons.DebugLog("NewConfig() Unmarshal apcHhost failed: ", err.Error())
		}
	}
	if len(hosts) == 0 {
		commons.DebugLog("NewConfig() load preferences apcHosts failed using defaults ")
		save, err := json.Marshal(defaultHosts)
		if err != nil {
			commons.DebugLog("NewConfig() Marshal saving apcHosts prefs failed: ", err.Error())
		}
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

	hubHostString := prefs.String(HubHostsPrefs)
	if hubHostString != "" {
		commons.DebugLog("NewConfig() load hub preferences succeeded ")
		err = json.Unmarshal([]byte(hubHostString), &hubHosts)
		if err != nil {
			commons.DebugLog("NewConfig() Unmarshal HubHosts failed: ", err.Error())
		}
	}
	if len(hubHosts) == 0 {
		commons.DebugLog("NewConfig() load hubHost preferences failed using defaults ")
		save, err := json.Marshal(defaultHubHosts)
		if err != nil {
			commons.DebugLog("NewConfig() Marshal saving hubHosts prefs failed: ", err.Error())
		}
		prefs.SetString(HubHostsPrefs, string(save))
		hubHosts = defaultHubHosts
	}

	cfg := &config{
		hosts: hosts,
		hubs:  hubHosts,
		prefs: prefs,
	}

	for _, h := range cfg.hosts {
		_ = cfg.VerifyHostConnection(h)
	}

	return cfg, err
}
func (c *config) ResetConfig() {
	c.prefs.SetString(HostsPrefs, "")
	c.prefs.SetString(HubHostsPrefs, "")
}
func (c *config) HostByName(hostName string) *entities.ApcHost {
	return c.hosts[hostName]
}
func (c *config) Hosts() []*entities.ApcHost {
	var r []*entities.ApcHost
	for _, v := range c.hosts {
		r = append(r, v)
	}
	return r
}
func (c *config) HubHosts() []*entities.HubHost {
	return c.hubs
}
func (c *config) Apply(h *entities.ApcHost) interfaces.Configuration {
	c.hosts[h.Name] = h
	_ = c.VerifyHostConnection(h)
	return c
}
func (c *config) ApplyHub(h *entities.HubHost) interfaces.Configuration {
	c.hubs = append(c.hubs, h)
	return c
}
func (c *config) AddHost(host *entities.ApcHost) {
	c.Apply(host).Save()
	commons.DebugLog("Config::AddHost() saved: .", host)
}
func (c *config) AddHubHost(host *entities.HubHost) {
	c.ApplyHub(host).Save()
	commons.DebugLog("Config::AddHubHost() saved: .", host)
}
func (c *config) Save() {
	save, err := json.Marshal(c.hosts)
	if err != nil {
		commons.DebugLog("Configuration::Save() marshal apcHosts failed: ", err.Error())
	} else {
		c.prefs.SetString(HostsPrefs, string(save))
	}
	save, err = json.Marshal(c.hubs)
	if err != nil {
		commons.DebugLog("Configuration::Save() marshal hubHosts failed: ", err.Error())
	} else {
		c.prefs.SetString(HubHostsPrefs, string(save))
	}
}
func (c *config) Remove(hostName string) {
	if hostName == "" {
		return
	}
	delete(c.hosts, hostName)
	c.Save()
}
func (c *config) HostKeys() []string {
	return commons.Keys(c.hosts)
}

// Shutdown compliance with Provider Interface
func (c *config) Shutdown() {
	commons.DebugLog("Config::Shutdown() called.")
}

// VerifyHostConnection compliance with Provider Interface
func (c *config) VerifyHostConnection(h *entities.ApcHost) error {
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
	time.Sleep(100 * time.Millisecond)
	h.State = commons.HostStatusOnline
	_ = conn.Close()
	return nil
}
