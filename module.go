package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Module is a component in the status bar.
type Module interface {
	// Name of  the module (this must be unique).
	Name() string
	// Update interval of the module.
	// If negative, the module is not updated on an interval basis - one would
	// use Signals to update it instead.
	UpdateInterval() time.Duration
	// Run module and return string to display.
	Exec() string
}

//============================
// External Modules
//============================

type ExternalModuleKind string

const (
	EMK_BIN ExternalModuleKind = "binary"
	EMK_SCR ExternalModuleKind = "script"
)

// ExternalModule is a Module implementation obtains its output by running a
// command.
// Depending in its kind, it attempts to execute a binary application (if
// SMK_BIN) or a shell script (if SMK_SCR) given by the command name with the
// given arguments.
// Note that the command should not return a non-zero error code as this is
// considered, and treated as an error (this will be logged so it can be
// relatively easy to detect).
// If this is out of your control (for example using a distributed binary), you
// can wrap it in a shell script and use SMK_SCR instead.
//
// If SMK_SCR, the script is executed using the shell pointed to by the
// constant SM_SHELL (by default, it is /bin/bash).
//
// If a postProcess function is provided, the output of the command is first
// processed through it before updating the status-bar.
type ExternalModule struct {
	name           string
	kind           ExternalModuleKind
	command        string   // name of the command to run
	args           []string // arguments to supply to command
	updateInterval time.Duration
	postProcess    func(out string) string
}

func (s *ExternalModule) Name() string                  { return s.name }
func (s *ExternalModule) UpdateInterval() time.Duration { return s.updateInterval }
func (s *ExternalModule) Exec() string {
	e := "EXT_MOD_EXEC_ERR(" + s.name + ")"
	var cmd *exec.Cmd
	switch s.kind {
	case EMK_BIN:
		cmd = exec.Command(s.command, s.args...)
	case EMK_SCR:
		args := []string{s.command}
		cmd = exec.Command("/bin/bash", append(args, s.args...)...)
	default:
		log(ERROR, "%s: unknown ExternalModuleKind '%s'", e, s.kind)
		return e
	}
	// run the command, post-process its output (if applicable) and return
	if outb, err := cmd.Output(); err != nil {
		log(ERROR, "%s: %s", e, err)
		return e
	} else if out := string(outb); s.postProcess != nil {
		return s.postProcess(out)
	} else {
		return out
	}
}

//============================
// Compiled Modules
//============================

type DateModule struct{}

func (d *DateModule) Name() string                  { return "date" }
func (d *DateModule) UpdateInterval() time.Duration { return time.Second }
func (d *DateModule) Exec() string {
	date, time := formattedDate()
	return fmt.Sprintf("%s %s %s %s", ICO_DATE, date, ICO_TIME, time)
}

func formattedDate() (d, t string) {
	now := time.Now()
	d = now.Format("Mon Jan 2, 2006")
	t = now.Format("15:04:05")
	return
}

//============================

// BatteryModule displays basic information about the battery pointed to by the
// device field (percentage capacity and whether it is plugged in).
type BatteryModule struct {
	device string // e.g. BAT1
}

func (b *BatteryModule) Name() string                  { return "battery" }
func (b *BatteryModule) UpdateInterval() time.Duration { return 30 * time.Second }
func (b *BatteryModule) Exec() string {
	e := "BAT_MOD_ERR"

	bat := fmt.Sprintf("/sys/class/power_supply/%s", b.device)
	capacityb, err := ioutil.ReadFile(bat + "/capacity")
	if err != nil {
		return e
	}
	statusb, err := ioutil.ReadFile(bat + "/status")
	if err != nil {
		return e
	}
	capacity := strings.TrimSpace(string(capacityb))
	status := strings.TrimSpace(string(statusb))
	return fmt.Sprintf("%s %s%%", batIcon(capacity, status), capacity)
}

const (
	STAT_FULL        = "Full"
	STAT_CHARGING    = "Charging"
	STAT_DISCHARGING = "Discharging"
)

// return an icon for the battery
func batIcon(cap, stat string) string {
	capacity, _ := strconv.ParseInt(cap, 10, 64)
	switch {
	case stat == STAT_CHARGING:
		switch {
		case 0 <= capacity && capacity <= 20:
			return ICO_BAT_C20
		case 20 < capacity && capacity <= 40:
			return ICO_BAT_C40
		case 40 < capacity && capacity <= 60:
			return ICO_BAT_C60
		case 60 < capacity && capacity <= 80:
			return ICO_BAT_C80
		default:
			return ICO_BAT_C90
		}
	case stat == STAT_FULL:
		return ICO_BAT_FULL
	case 0 <= capacity && capacity <= 5:
		return ICO_BAT_00
	case 5 < capacity && capacity <= 10:
		return ICO_BAT_10
	case 10 < capacity && capacity <= 20:
		return ICO_BAT_20
	case 20 < capacity && capacity <= 30:
		return ICO_BAT_30
	case 30 < capacity && capacity <= 40:
		return ICO_BAT_40
	case 40 < capacity && capacity <= 50:
		return ICO_BAT_50
	case 50 < capacity && capacity <= 60:
		return ICO_BAT_60
	case 60 < capacity && capacity <= 70:
		return ICO_BAT_70
	case 70 < capacity && capacity <= 80:
		return ICO_BAT_80
	case 80 < capacity && capacity <= 90:
		return ICO_BAT_90
	default:
		return ICO_BAT_100
	}
}

//============================

var MasterVolumeModule ExternalModule = ExternalModule{
	name:           "volume",
	command:        "vol",
	kind:           EMK_SCR,
	updateInterval: -1,
	args:           []string{"get"},
	postProcess: func(out string) string {
		e := "EXT_MOD_MVOL_ERR"

		out = strings.TrimSpace(out)
		if out == "muted" {
			return ICO_VOL_MUTE
		} else if vol, err := strconv.ParseInt(out[:len(out)-1], 10, 32); err != nil {
			log(ERROR, "%s: %s", e, err)
			return e
		} else if 0 <= vol && vol < 50 {
			return ICO_VOL_DOWN + " " + out
		} else {
			return ICO_VOL_UP + " " + out
		}
	},
}

//============================

var RAMUsageModule ExternalModule = ExternalModule{
	name:           "sys_ram_usage",
	command:        "sb_ram_usage",
	kind:           EMK_SCR,
	updateInterval: 5 * time.Second,
	postProcess: func(out string) string {
		e := "EXT_MOD_RAM_ERR"

		parts := strings.Split(out, " ")
		if len(parts) != 2 {
			log(ERROR, "%s: malformed script output '%s'", e, out)
			return e
		}
		useds := strings.TrimSpace(parts[0])
		totals := strings.TrimSpace(parts[1])
		if used, err := strconv.ParseInt(useds, 10, 64); err != nil {
			log(ERROR, "%s: invalid used size '%s'", e, useds)
			return e
		} else if total, err := strconv.ParseInt(totals, 10, 64); err != nil {
			log(ERROR, "%s: invalid total size '%s'", e, totals)
			return e
		} else {
			return fmt.Sprintf("%s %d%%", ICO_RES_RAM, 100*used/total)
		}
	},
}

var CPUUsageModule ExternalModule = ExternalModule{
	name:           "sys_cpu_usage",
	command:        "sb_cpu_usage",
	kind:           EMK_SCR,
	updateInterval: 5 * time.Second,
	postProcess: func(out string) string {
		return ICO_RES_CPU + " " + strings.TrimSpace(out)
	},
}
