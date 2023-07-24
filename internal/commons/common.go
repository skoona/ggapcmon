/**
 * commons
 * is the collector of common utilities used
 */

package commons

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	AppIcon             = "apcupsd"
	PreferencesIcon     = "preferences"
	HostStatusUnknown   = "unknown"
	HostStatusOnBattery = "onbattery"
	HostStatusCharging  = "charging"
	HostStatusOnline    = "online"
)

// ShutdownSignals alternate panic() implementation, causes an orderly shutdown
var ShutdownSignals chan os.Signal
var DebugLoggingEnabled = ("true" == os.Getenv("GAPC_DEBUG")) // "true" / "false"
var logs = log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)

func DebugLog(args ...any) {
	if DebugLoggingEnabled {
		_ = logs.Output(2, fmt.Sprint(args...))
	}
}

// Keys returns the keys of the map m.
// The keys will be an indeterminate order.
// alternate reflect based: reflect.ValueOf(m).MapKeys()
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// ChangeTimeFormat converts APC timestamp to something more human readable
// time.RFC1123, time.RFC3339 are good choices
// returns local time version of value
func ChangeTimeFormat(timeString string, format string) string {
	if format == "" {
		format = time.RFC1123
	}
	if timeString == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", strings.TrimSpace(timeString))
	if err != nil {
		DebugLog("ApcService::ChangeTimeFormat() Time Parse Error, src: ", timeString, ", err: ", err.Error())
	}
	return t.Format(format)
}

// ShiftSlice drops index 0 and append newData to any type of slice
func ShiftSlice[K comparable](newData K, slice []K) []K {
	idx := 0
	if len(slice) == 0 {
		return append(slice, newData)
	}
	shorter := append(slice[:idx], slice[idx+1:]...)
	shorter = append(shorter, newData)
	return shorter
}
