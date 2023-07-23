package providers

import (
	"context"
	"encoding/json"
	"fyne.io/fyne/v2"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"log"
	"net"
	"time"
)

const (
	HostLocal     = "127.0.0.1:3551"
	HostLocalName = "Local"
	HostVServ     = "10.100.1.3:3551"
	HostVServName = "VServ"
	HostPve       = "10.100.1.4:3551"
	HostPveName   = "PVE"
	HostsPrefs    = "ApcHost"
)

type config struct {
	hosts map[string]*entities.ApcHost
	log   *log.Logger
	prefs fyne.Preferences
}

var _ interfaces.Configuration = (*config)(nil)
var _ interfaces.Provider = (*config)(nil)

func NewConfig(prefs fyne.Preferences, log *log.Logger) (interfaces.Configuration, error) {
	var err error
	var hosts map[string]*entities.ApcHost

	defaultHosts := map[string]*entities.ApcHost{
		// graph-30 = 15 hours @ 15 network-sec
		HostLocalName: entities.NewApcHost(HostLocalName, HostLocal, 10, 5, true, true),
		HostVServName: entities.NewApcHost(HostVServName, HostVServ, 10, 5, true, true),
		HostPveName:   entities.NewApcHost(HostPveName, HostPve, 10, 5, true, true),
	}

	hostString := prefs.String(HostsPrefs)
	if hostString != "" && len(hostString) > 16 {
		log.Println("NewConfig() load preferences succeeded ")
		err = json.Unmarshal([]byte(hostString), &hosts)
		if err != nil {
			log.Println("NewConfig() Unmarshal failed: ", err.Error())
		}
	}
	if len(hosts) == 0 {
		log.Println("NewConfig() load preferences failed using defaults ")
		save, err := json.Marshal(defaultHosts)
		if err != nil {
			log.Println("NewConfig() Marshal saving prefs failed: ", err.Error())
		}
		prefs.SetString(HostsPrefs, string(save))
		hosts = defaultHosts
	}

	cfg := &config{
		hosts: hosts,
		log:   log,
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
func (c *config) Apply(h *entities.ApcHost) interfaces.Configuration {
	c.hosts[h.Name] = h
	_ = c.VerifyHostConnection(h)
	return c
}
func (c *config) AddHost(host *entities.ApcHost) {
	c.Apply(host).Save()
	c.log.Println("Config::AddHost() saved: .", host)
}
func (c *config) Save() {
	save, err := json.Marshal(c.hosts)
	if err != nil {
		log.Println("Configuration::Save() marshal failed: ", err.Error())
	} else {
		c.prefs.SetString(HostsPrefs, string(save))
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
	c.log.Println("Config::Shutdown() called.")
}

// VerifyHostConnection compliance with Provider Interface
func (c *config) VerifyHostConnection(h *entities.ApcHost) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", h.IpAddress)
	if err != nil {
		log.Println("connect() dial Error: ", err.Error(), ", Ip: ", h.IpAddress, ", context: ", ctx.Err())
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
