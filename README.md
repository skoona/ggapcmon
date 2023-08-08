# ggApcMon

A GO/Fyne refactor of gapcmon from sourceforge, originally written in C/Gtk2. 

Monitors and charts the metrics of a APC UPS which use the [APCUPSD](http://www.apcupsd.org) software package.




### Project Layout

Enable debug logging via environment var: `export GAPC_DEBUG="true"`


```text
├── FyneApp.toml
├── LICENSE
├── Makefile
├── README.md
├── bin
│   └── ggapcmon
├── docs
│   ├── statusFormats.md
│   └── summary-c.md
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── handlers
│   │   │   └── ui
│   │   │       ├── detailed.go
│   │   │       ├── glossary.go
│   │   │       ├── menus.go
│   │   │       ├── monitor.go
│   │   │       ├── overview.go
│   │   │       ├── preferences.go
│   │   │       └── viewprovider.go
│   │   └── repository
│   │       ├── apcprovider.go
│   │       └── config.go
│   ├── commons
│   │   ├── common.go
│   │   ├── imageResources.go
│   │   ├── images.go
│   │   └── resources
│   │       ├── apcupsd.png
│   │       ├── charging.png
│   │       ├── gapc_prefs.png
│   │       ├── onbatt.png
│   │       ├── online.png
│   │       └── unplugged.png
│   └── core
│       ├── domain
│       │   ├── apchosts.go
│       │   ├── channelTuple.go
│       │   ├── graphaverage.go
│       │   └── upsstatusvaluebindings.go
│       ├── ports
│       │   ├── apcprovider.go
│       │   ├── configuration.go
│       │   ├── graphpointsmoothing.go
│       │   ├── provider.go
│       │   └── service.go
│       └── services
│           └── service.go
├── main.go
└── skoona.png
```

### Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request


### LICENSE
The application is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
