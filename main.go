package main

import "time"

func main() {
	bar := StatusBar{
		Sep: " | ",
		Modules: []Module{
			&DateModule{},
			&BatteryModule{"BAT1"},
		},
		UpdateBatchSize: 3,
		UpdateMinTime:   time.Second,
	}
	bar.Exec()
}
