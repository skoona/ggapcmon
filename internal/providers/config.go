package providers

import (
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"log"
	"strings"
	"time"
)

const (
	HostVServ     = "10.100.1.3:3551"
	HostVServName = "VServ"
	HostPve       = "10.100.1.4:3551"
	HostPveName   = "PVE"
	HostsPrefs    = "ApcHost"
)

type config struct {
	hosts map[string]entities.ApcHost
	log   *log.Logger
	prefs fyne.Preferences
}

var _ interfaces.Configuration = (*config)(nil)
var _ interfaces.Provider = (*config)(nil)

func NewConfig(prefs fyne.Preferences, log *log.Logger) (interfaces.Configuration, error) {
	var err error
	var hosts map[string]entities.ApcHost

	defaultHosts := map[string]entities.ApcHost{
		// graph-30 = 15 hours @ 15 network-sec
		HostVServName: entities.ApcHost{IpAddress: HostVServ, Name: HostVServName, NetworkSamplePeriod: 15, GraphingSamplePeriod: 5, Enabled: true, TrayIcon: true},
		HostPveName:   entities.ApcHost{IpAddress: HostPve, Name: HostPveName, NetworkSamplePeriod: 15, GraphingSamplePeriod: 5, Enabled: true, TrayIcon: true},
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

	return cfg, err
}
func (c *config) ResetConfig() {
	c.prefs.SetString(HostsPrefs, "")
}
func (c *config) HostByName(hostName string) entities.ApcHost {
	return c.hosts[hostName]
}
func (c *config) Hosts() []entities.ApcHost {
	var r []entities.ApcHost
	for _, v := range c.hosts {
		r = append(r, v)
	}
	return r
}
func (c *config) Apply(h entities.ApcHost) entities.ApcHost {
	c.hosts[h.Name] = h
	return c.hosts[h.Name]
}
func (c *config) Save(hosts []entities.ApcHost) error {
	if hosts == nil {
		return errors.New("Configuration::Save() host parameter cannot be nil")
	}
	for _, host := range hosts {
		c.hosts[host.Name] = host
	}
	save, err := json.Marshal(c.hosts)
	if err != nil {
		log.Println("Configuration::Save() marshal failed: ", err.Error())
	} else {
		c.prefs.SetString(HostsPrefs, string(save))
	}

	return err
}
func (c *config) Update(name, ip string, netperiod, graphperiod time.Duration, tray, enable bool) entities.ApcHost {
	host := c.hosts[name]
	host.Name = strings.Clone(name)
	host.IpAddress = strings.Clone(ip)
	host.NetworkSamplePeriod = netperiod
	host.GraphingSamplePeriod = graphperiod
	host.TrayIcon = tray
	host.Enabled = enable

	return c.hosts[name]
}
func (c *config) Remove(hostName string) {
	if hostName == "" {
		return
	}
	delete(c.hosts, hostName)
}
func (c *config) HostKeys() []string {
	return commons.Keys(c.hosts)
}

// Shutdown closes all go routine
func (c *config) Shutdown() {
	c.log.Println("Config::Shutdown() called.")
}
