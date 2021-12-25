package main

import (
	"flag"
	"time"
)

var production = flag.Bool("prod", true, "set to 'false' to log excessively")
var SIG_SOCK = flag.String("sock", "", "name of socket to create for status bar signalling")

const (
	// shell used to execute scripted modules
	SM_SHELL = "/bin/bash"
)

// Icons
// These are nerd-font icons, however even emojis can be used - go nuts!
// Reference: https://www.nerdfonts.com/cheat-sheet
const (
	// DateModule
	ICO_DATE = ""
	ICO_TIME = ""
	// BatteryModule
	ICO_BAT_00   = ""
	ICO_BAT_10   = ""
	ICO_BAT_20   = ""
	ICO_BAT_30   = ""
	ICO_BAT_40   = ""
	ICO_BAT_50   = ""
	ICO_BAT_60   = ""
	ICO_BAT_70   = ""
	ICO_BAT_80   = ""
	ICO_BAT_90   = ""
	ICO_BAT_100  = ""
	ICO_BAT_C20  = ""
	ICO_BAT_C40  = ""
	ICO_BAT_C60  = ""
	ICO_BAT_C80  = ""
	ICO_BAT_C90  = ""
	ICO_BAT_FULL = "(F)"
	// VolumeModule
	ICO_VOL_MUTE = "婢"
	ICO_VOL_DOWN = "奔"
	ICO_VOL_UP   = "墳"
	// ResourcesModule
	ICO_RES_RAM = ""
	ICO_RES_CPU = ""
)

// This is the bar which will be rendered.
var statusBar StatusBar = StatusBar{
	Sep: " ",
	Modules: []Module{
		&CPUUsageModule,
		&RAMUsageModule,
		&DateModule{},
		&MasterVolumeModule,
		&BatteryModule{"BAT1"},
	},
	// update the bar whenever a module update is available
	UpdateBatchSize: 1,
	UpdateMinTime:   time.Second,
}
