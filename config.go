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
// These are font-awesome icons, however even emojis can be used - go nuts!
const (
	// DateModule
	ICO_DATE = ""
	ICO_TIME = ""
	// BatteryModule
	ICO_BAT_POW  = ""
	ICO_BAT_FULL = "(F)"
	ICO_BAT_0Q   = ""
	ICO_BAT_1Q   = ""
	ICO_BAT_2Q   = ""
	ICO_BAT_3Q   = ""
	ICO_BAT_4Q   = ""
	// VolumeModule
	ICO_VOL_MUTE = ""
	ICO_VOL_DOWN = ""
	ICO_VOL_UP   = ""
	// ResourcesModule
	ICO_RES_RAM = ""
	ICO_RES_CPU = ""
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
