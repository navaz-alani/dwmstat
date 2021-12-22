package main

import "time"

const (
	// shell used to execute scripted modules
	SM_SHELL = "/bin/bash"
)

// Icons
const (
	// DateModule
	ICO_DATE = " "
	ICO_TIME = " "
	// BatteryModule
	ICO_BAT_POW  = ""
	ICO_BAT_FULL = " (F)"
	ICO_BAT_0Q   = " "
	ICO_BAT_1Q   = " "
	ICO_BAT_2Q   = " "
	ICO_BAT_3Q   = " "
	ICO_BAT_4Q   = " "
	// VolumeModule
	ICO_VOL_MUTE = " "
	ICO_VOL_DOWN = " "
	ICO_VOL_UP   = " "
)

// This is the bar which will be rendered.
var statusBar StatusBar = StatusBar{
	Sep: " | ",
	Modules: []Module{
		&DateModule{},
		&BatteryModule{"BAT1"},
		&MasterVolumeModule,
	},
	// update the bar whenever a module update is available
	UpdateBatchSize: 1,
	UpdateMinTime:   time.Second,
}
