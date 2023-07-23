package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (v *viewProvider) GlossaryPage() *fyne.Container {
	gtext := `
# GAPCMON 

A monitor for UPS's under the management of APCUPSD. 

When active, gapcmon provides three visual objects to interact with. 
First is the main control.panel where monitors are defined, enabled, and listed when 
active. Second are notification area icons that manage the visibility of each window or 
panel. The third is an information window showing historical and current details of 
a UPS being monitored. 


## CONTROL PANEL WINDOW PAGES
***ACTIVE MONITORS PAGE***

A short list of the monitors that are currently enabled.  The list shows 
each monitor's current icon, its status, and a brief summary of its key metrics.  
Double-clicking a row causes the information window of that monitor to be presented.

### PREFERENCES PAGE
-for the monitors

#### Enable:
Causes the monitor to immediately run, create an info-window, 
and add an entry in the icon list window.

#### Use Trayicon:
Adds a notification area icon which toggles 
the visibility of the monitor info-window when clicked.

#### Network refresh:
The number of seconds between collections of status and event 
data from the network.

#### Graph refresh:
Multiplied by the network refresh value to determine 
the total number of seconds between graph data collections.

#### Hostname or IP Address:
The hostname or address where an apcupsd NIS interface 
is running.

#### Port:
The NIS access port on the APCUPSD host; defaults to 3551.

#### Add | Remove Buttons:
Buttons to add or remove a monitor entry from the 
list of monitors above.  Add, adds with defaulted values at end of list.  Remove, removes 
the currently selected row.

> _for the control.panel_

#### Use Trayicon:
Adds a notification area icon which toggles 
the visibility of the control panel when clicked.  **Note:** _all tray icons 
contain a popup menu with the choices of 'JumpTo' interactive window, and 'Quit' 
which either hides the window or destroys it in the case of the control panel. 
Additionally, when use_trayicon is selected, the title of the window is removed from 
the desktop's windowlist or taskbar._


#### GRAPH PROPERTIES PAGE
Allows you to specify the colors to be used for each of the five data series on the 
Historical Summary graph page of each monitor.   General window background colors can 
also be specified.


### GLOSSARY PAGE

This page of introductory text.

#### ABOUT PAGE
More standard vanity, and my e-mail ID in case something breaks.

## MONITOR INFORMATION WINDOW PAGES

#### HISTORICAL SUMMARY PAGE
A graph showing the last 40 samples of five key data points, scaled to represent all 
points as a percentage of that value's normal range.  A data point's value can be 
viewed by moving the mouse over any desired point, a tooltip will appear 
showing the color and value of all points at that interval.  Data points are collected 
periodically, based on the product of graph_refresh times network_refresh in seconds. 
These tooltips can be enabled or disabled by clicking anywhere on the graph once.

#### DETAILED INFORMATION PAGE
A more in-depth view of the monitored UPS's environmental values.  Software, product, 
and operational values are available and updated every 'network_refresh' seconds.

#### POWER EVENTS PAGE
A log of all power events recorded by APCUPSD on the server.

#### FULL UPS STATUS PAGE
A listing of the output from apcaccess showing the actual state as reported 
by the UPS.
`

	rtext := widget.NewRichTextFromMarkdown("")
	rtext.ParseMarkdown(gtext)
	rtext.Wrapping = fyne.TextWrapWord
	return container.NewMax(
		container.NewVScroll(
			rtext,
		),
	)
}
