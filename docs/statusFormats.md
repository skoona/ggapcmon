# Status Formats are different for Master/Slave
```text

[DEBUG] 09:43:28.671942 main.go:73: {VServ}(0)[status]  ==> case "APC      ":// 001,021,0577
[DEBUG] 09:43:28.671969 main.go:73: {VServ}(1)[status]  ==> case "DATE     ":// Thu, 20 Jul 2023 09:43:04 EDT
[DEBUG] 09:43:28.671973 main.go:73: {VServ}(2)[status]  ==> case "HOSTNAME ":// vserv
[DEBUG] 09:43:28.671975 main.go:73: {VServ}(3)[status]  ==> case "VERSION  ":// 3.14.14 (31 May 2016) debian
[DEBUG] 09:43:28.671977 main.go:73: {VServ}(4)[status]  ==> case "UPSNAME  ":// pve
[DEBUG] 09:43:28.671979 main.go:73: {VServ}(5)[status]  ==> case "CABLE    ":// Ethernet Link
[DEBUG] 09:43:28.671982 main.go:73: {VServ}(6)[status]  ==> case "DRIVER   ":// NETWORK UPS Driver
[DEBUG] 09:43:28.671984 main.go:73: {VServ}(7)[status]  ==> case "UPSMODE  ":// Stand Alone
[DEBUG] 09:43:28.671986 main.go:73: {VServ}(8)[status]  ==> case "STARTTIME":// Sun, 02 Apr 2023 14:17:09 EDT
[DEBUG] 09:43:28.671989 main.go:73: {VServ}(9)[status]  ==> case "MASTERUPD":// Thu, 20 Jul 2023 09:43:04 EDT
[DEBUG] 09:43:28.671994 main.go:73: {VServ}(10)[status] ==> case "MASTER   ":// 10.100.1.4:3551
[DEBUG] 09:43:28.671996 main.go:73: {VServ}(11)[status] ==> case "STATUS   ":// ONLINE SLAVE
[DEBUG] 09:43:28.672006 main.go:73: {VServ}(12)[status] ==> case "MBATTCHG ":// 10 Percent
[DEBUG] 09:43:28.672009 main.go:73: {VServ}(13)[status] ==> case "MINTIMEL ":// 5 Minutes
[DEBUG] 09:43:28.672012 main.go:73: {VServ}(14)[status] ==> case "MAXTIME  ":// 0 Seconds
[DEBUG] 09:43:28.672014 main.go:73: {VServ}(15)[status] ==> case "NUMXFERS ":// 1
[DEBUG] 09:43:28.672016 main.go:73: {VServ}(16)[status] ==> case "XONBATT  ":// Mon, 17 Jul 2023 12:02:00 EDT
[DEBUG] 09:43:28.672019 main.go:73: {VServ}(17)[status] ==> case "TONBATT  ":// 0 Seconds
[DEBUG] 09:43:28.672021 main.go:73: {VServ}(18)[status] ==> case "CUMONBATT":// 2 Seconds
[DEBUG] 09:43:28.672025 main.go:73: {VServ}(19)[status] ==> case "XOFFBATT ":// Mon, 17 Jul 2023 12:02:02 EDT
[DEBUG] 09:43:28.672028 main.go:73: {VServ}(20)[status] ==> case "STATFLAG ":// 0x05000408
[DEBUG] 09:43:28.672030 main.go:73: {VServ}(21)[status] ==> case "END APC  ":// Thu, 20 Jul 2023 09:43:27 EDT
=============
[DEBUG] 09:43:30.214919 main.go:65: {PVE}(0)[status]  ==> case "APC      ":// 001,045,1108
[DEBUG] 09:43:30.214936 main.go:65: {PVE}(1)[status]  ==> case "DATE     ":// Thu, 20 Jul 2023 09:43:05 EDT
[DEBUG] 09:43:30.214940 main.go:65: {PVE}(2)[status]  ==> case "HOSTNAME ":// pve
[DEBUG] 09:43:30.214943 main.go:65: {PVE}(3)[status]  ==> case "VERSION  ":// 3.14.14 (31 May 2016) debian
[DEBUG] 09:43:30.214945 main.go:65: {PVE}(4)[status]  ==> case "UPSNAME  ":// pve
[DEBUG] 09:43:30.214947 main.go:65: {PVE}(5)[status]  ==> case "CABLE    ":// USB Cable
[DEBUG] 09:43:30.214949 main.go:65: {PVE}(6)[status]  ==> case "DRIVER   ":// USB UPS Driver
[DEBUG] 09:43:30.214951 main.go:65: {PVE}(7)[status]  ==> case "UPSMODE  ":// Stand Alone
[DEBUG] 09:43:30.214954 main.go:65: {PVE}(8)[status]  ==> case "STARTTIME":// Sun, 02 Apr 2023 14:15:31 EDT
[DEBUG] 09:43:30.214956 main.go:65: {PVE}(9)[status]  ==> case "MODEL    ":// Smart-UPS 1500
[DEBUG] 09:43:30.214958 main.go:65: {PVE}(10)[status] ==> case "STATUS   ":// ONLINE
[DEBUG] 09:43:30.214960 main.go:65: {PVE}(11)[status] ==> case "LINEV    ":// 125.2 Volts
[DEBUG] 09:43:30.214962 main.go:65: {PVE}(12)[status] ==> case "LOADPCT  ":// 39.6 Percent
[DEBUG] 09:43:30.214978 main.go:65: {PVE}(13)[status] ==> case "BCHARGE  ":// 100.0 Percent
[DEBUG] 09:43:30.214981 main.go:65: {PVE}(14)[status] ==> case "TIMELEFT ":// 29.0 Minutes
[DEBUG] 09:43:30.214983 main.go:65: {PVE}(15)[status] ==> case "MBATTCHG ":// 5 Percent
[DEBUG] 09:43:30.214985 main.go:65: {PVE}(16)[status] ==> case "MINTIMEL ":// 3 Minutes
[DEBUG] 09:43:30.214987 main.go:65: {PVE}(17)[status] ==> case "MAXTIME  ":// 0 Seconds
[DEBUG] 09:43:30.214990 main.go:65: {PVE}(18)[status] ==> case "OUTPUTV  ":// 125.2 Volts
[DEBUG] 09:43:30.214995 main.go:65: {PVE}(19)[status] ==> case "SENSE    ":// High
[DEBUG] 09:43:30.214998 main.go:65: {PVE}(20)[status] ==> case "DWAKE    ":// -1 Seconds
[DEBUG] 09:43:30.215000 main.go:65: {PVE}(21)[status] ==> case "DSHUTD   ":// 90 Seconds
[DEBUG] 09:43:30.215002 main.go:65: {PVE}(22)[status] ==> case "LOTRANS  ":// 106.0 Volts
[DEBUG] 09:43:30.215004 main.go:65: {PVE}(23)[status] ==> case "HITRANS  ":// 127.0 Volts
[DEBUG] 09:43:30.215006 main.go:65: {PVE}(24)[status] ==> case "RETPCT   ":// 0.0 Percent
[DEBUG] 09:43:30.215136 main.go:65: {PVE}(25)[status] ==> case "ITEMP    ":// 29.7 C
[DEBUG] 09:43:30.215143 main.go:65: {PVE}(26)[status] ==> case "ALARMDEL ":// 30 Seconds
[DEBUG] 09:43:30.215146 main.go:65: {PVE}(27)[status] ==> case "BATTV    ":// 27.4 Volts
[DEBUG] 09:43:30.215148 main.go:65: {PVE}(28)[status] ==> case "LINEFREQ ":// 60.0 Hz
[DEBUG] 09:43:30.215150 main.go:65: {PVE}(29)[status] ==> case "LASTXFER ":// Automatic or explicit self test
[DEBUG] 09:43:30.215154 main.go:65: {PVE}(30)[status] ==> case "NUMXFERS ":// 12
[DEBUG] 09:43:30.215157 main.go:65: {PVE}(31)[status] ==> case "XONBATT  ":// Mon, 17 Jul 2023 12:01:53 EDT
[DEBUG] 09:43:30.215159 main.go:65: {PVE}(32)[status] ==> case "TONBATT  ":// 0 Seconds
[DEBUG] 09:43:30.215161 main.go:65: {PVE}(33)[status] ==> case "CUMONBATT":// 71 Seconds
[DEBUG] 09:43:30.215163 main.go:65: {PVE}(34)[status] ==> case "XOFFBATT ":// Mon, 17 Jul 2023 12:02:00 EDT
[DEBUG] 09:43:30.215165 main.go:65: {PVE}(35)[status] ==> case "LASTSTEST":// Mon, 17 Jul 2023 12:01:53 EDT
[DEBUG] 09:43:30.215168 main.go:65: {PVE}(36)[status] ==> case "SELFTEST ":// NO
[DEBUG] 09:43:30.215215 main.go:65: {PVE}(37)[status] ==> case "STESTI   ":// 14 days
[DEBUG] 09:43:30.215219 main.go:65: {PVE}(38)[status] ==> case "STATFLAG ":// 0x05000008
[DEBUG] 09:43:30.215221 main.go:65: {PVE}(39)[status] ==> case "MANDATE  ":// 2002-06-14
[DEBUG] 09:43:30.215223 main.go:65: {PVE}(40)[status] ==> case "SERIALNO ":// AS0224131174
[DEBUG] 09:43:30.215225 main.go:65: {PVE}(41)[status] ==> case "BATTDATE ":// 2021-09-06
[DEBUG] 09:43:30.215227 main.go:65: {PVE}(42)[status] ==> case "NOMOUTV  ":// 120 Volts
[DEBUG] 09:43:30.215229 main.go:65: {PVE}(43)[status] ==> case "NOMBATTV ":// 24.0 Volts
[DEBUG] 09:43:30.215232 main.go:65: {PVE}(44)[status] ==> case "FIRMWARE ":// 601.3.D USB FW:1.3
[DEBUG] 09:43:30.215234 main.go:65: {PVE}(45)[status] ==> case "END APC  ":// Thu, 20 Jul 2023 09:43:27 EDT
```
