package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// Module is a component in the status bar.
type Module interface {
	// Name of  the module.
	Name() string
	// Update interval
	UpdateInterval() time.Duration
	// Update signal used to trigger module update (-1 if none).
	UpdateSig() int
	// Run module and return string to display.
	Exec() string
}

//============================
// Scripted Modules
//============================

// ScriptedModule is a Module implementation obtains its output
type ScriptedModule struct {
	name           string
	command        string   // name of the command to run
	args           []string // arguments to supply to command
	updateInterval time.Duration
	updateSig      int
}

func (s *ScriptedModule) Name() string                  { return s.name }
func (s *ScriptedModule) UpdateInterval() time.Duration { return s.updateInterval }
func (s *ScriptedModule) UpdateSig() int                { return s.updateSig }

//============================
// Compiled Modules
//============================

type DateModule struct{}

func (d *DateModule) Name() string                  { return "date" }
func (d *DateModule) UpdateInterval() time.Duration { return time.Second }
func (d *DateModule) UpdateSig() int                { return -1 }
func (d *DateModule) Exec() string {
	date, time := formattedDate()
	return fmt.Sprintf("%s %s %s %s", ICO_DATE, date, ICO_TIME, time)
}

const (
	ICO_DATE = " "
	ICO_TIME = " "
)

func formattedDate() (d, t string) {
	now := time.Now()
	d = now.Format("Mon Jan 2 2006")
	t = now.Format("15:04:05")
	return
}

//============================

type BatteryModule struct {
	device string // e.g. BAT1
}

func (b *BatteryModule) Name() string                  { return "battery" }
func (b *BatteryModule) UpdateInterval() time.Duration { return 30 * time.Second }
func (b *BatteryModule) UpdateSig() int                { return -1 }
func (b *BatteryModule) Exec() string {
	e := "BAT_MOD_ERR"

	bat := fmt.Sprintf("/sys/class/power_supply/%s", b.device)
	capacity, err := ioutil.ReadFile(bat + "/capacity")
	if err != nil {
		return e
	}
	status, err := ioutil.ReadFile(bat + "/status")
	if err != nil {
		return e
	}
	return fmt.Sprintf("%s %s%%", batIcon(string(capacity), string(status)), string(capacity))
}

const (
	STAT_FULL        = "Full"
	STAT_CHARGING    = "Charging"
	STAT_DISCHARGING = "Discharging"

	ICO_BAT_POW  = ""
	ICO_BAT_FULL = " (F)"
	ICO_BAT_0Q   = " "
	ICO_BAT_1Q   = " "
	ICO_BAT_2Q   = " "
	ICO_BAT_3Q   = " "
	ICO_BAT_4Q   = " "
)

// return an icon for the battery
func batIcon(cap, stat string) string {
	capacity, _ := strconv.ParseInt(strings.TrimSpace(cap), 10, 64)
	println(capacity)
	stat = strings.TrimSpace(stat)
	switch {
	case stat == STAT_CHARGING:
		return ICO_BAT_POW
	case stat == STAT_FULL:
		return ICO_BAT_FULL
	case capacity == 0:
		return ICO_BAT_0Q
	case 0 < capacity && capacity <= 25:
		return ICO_BAT_1Q
	case 25 < capacity && capacity <= 50:
		return ICO_BAT_2Q
	case 50 < capacity && capacity <= 75:
		return ICO_BAT_3Q
	default:
		return ICO_BAT_4Q
	}
}
