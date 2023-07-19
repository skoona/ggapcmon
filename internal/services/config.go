package services

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"log"
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
}

var _ interfaces.Configuration = (*config)(nil)

func NewConfig(prefs fyne.Preferences, log *log.Logger) (interfaces.Configuration, error) {
	var err error
	hosts := map[string]entities.ApcHost{}

	defaultHosts := map[string]entities.ApcHost{
		HostVServ:   entities.ApcHost{IpAddress: HostVServ, Name: HostVServName, SecondsPerSample: 33},
		HostPveName: entities.ApcHost{IpAddress: HostPve, Name: HostPveName, SecondsPerSample: 37},
	}

	hostString := prefs.String(HostsPrefs)
	if hostString != "" {
		log.Println("NewConfig() load preferences succeeded ")
		err = json.Unmarshal([]byte(hostString), &hosts)
		if err != nil {
			log.Println("NewConfig() Unmarshal failed: ", err.Error())
		}
	} else {
		log.Println("NewConfig() load preferences failed using defaults ")
		save, err := json.Marshal(hosts)
		if err != nil {
			log.Println("NewConfig() Marshal saving prefs failed: ", err.Error())
		}
		prefs.SetString(HostsPrefs, string(save))
		hosts = defaultHosts
	}

	cfg := &config{
		hosts: hosts,
		log:   log,
	}

	return cfg, err
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
