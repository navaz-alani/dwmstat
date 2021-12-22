package main

import (
	"fmt"
	"os/exec"
	"time"
)

func xsetroot(s string) error {
	cmd := exec.Command("xsetroot", "-name", s)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xsetroot('%s'): %s", s, err)
	}
	return nil
}

type StatusBar struct {
	Sep             string   // module separator
	Modules         []Module // list of modules
	UpdateBatchSize int      // number of updates to batch
	// time to wait before updating bar if batch is not big enough (see
	// StatusBar.barUpdater)
	UpdateMinTime time.Duration
}

func (b *StatusBar) Run() {
	signalStream := make(chan Signal, 5)
	modEvStreams := b.init()
	go listenSignals(signalStream)
	b.handleSignals(modEvStreams, signalStream)
}

// Event is a kind of signal sent to a worker.
type Event string

const (
	EV_STOP = "STOP" // skip current update internal and stop
	EV_EXEC = "EXEC" // don't wait for interval to finish, run now
)

type UpdateCallback func(out string)

// moduleUpdater is the function which periodically executes a module and
// submits its output.
// It can also receive and process events through the evStream, so if a module
// is not to be periodically updated, it can be updated through a signal
// instead.
func moduleUpdater(mod Module, evStream chan Event, cb UpdateCallback) {
	if mod.UpdateInterval() > 0 {
		select {
		case ev := <-evStream:
			switch ev {
			case EV_EXEC: // run without waiting for interval to complete
				if !*production {
					log(INFO, "EV_EXEC - updating module '%s'", mod.Name())
				}
				cb(mod.Exec())
			case EV_STOP:
				return
			}
		case <-time.After(mod.UpdateInterval()):
			if !*production {
				log(INFO, "INT_EXEC - updating module '%s'", mod.Name())
			}
			cb(mod.Exec())
		}
	} else {
		ev := <-evStream
		switch ev {
		case EV_EXEC:
			if !*production {
				log(INFO, "EV_EXEC - updating module '%s'", mod.Name())
			}
			cb(mod.Exec())
		case EV_STOP:
			return
		}
	}
	go moduleUpdater(mod, evStream, cb) // respawn
}

// moduleUpdateNotification is a message sent when a Module has updated.
type moduleUpdateNotification struct {
	idx int
	out string
}

type moduleEventStream chan Event

// mapping of module name to Module
func (b *StatusBar) init() (modEvStreams map[string]moduleEventStream) {
	modules := make(map[string]Module)
	// module event streams -> channels through which modules receive events
	modEvStreams = make(map[string]moduleEventStream)

	// collect all modules
	for _, m := range b.Modules {
		name := m.Name()
		if _, ok := modules[name]; ok {
			log(WARN, "duplicate module name '%s'", name)
		}
		modules[name] = m
		modEvStreams[name] = make(moduleEventStream)
	}

	// list of module outputs
	outs := make([]string, len(b.Modules))

	// initially run all modules
	for i, m := range b.Modules {
		outs[i] = m.Exec()
	}
	// channel over which module update routines notify updates
	var updateNotifyStream = make(chan moduleUpdateNotification)
	// dispatch a module update routine for each module
	for i, mod := range b.Modules {
		i := i
		go moduleUpdater(
			mod,
			modEvStreams[mod.Name()],
			// when the module updates, send a notification to the bar updater
			func(o string) { updateNotifyStream <- moduleUpdateNotification{i, o} },
		)
	}
	// dispatch routine that collects module updates and updates the status bar
	go b.barUpdater(outs, updateNotifyStream)
	return modEvStreams
}

// handleSignals handles submitted signals.
func (b *StatusBar) handleSignals(
	modEvStreams map[string]moduleEventStream,
	signalStream chan Signal,
) {
	// listen for any signals submitted through the socket
	for {
		select {
		case sig := <-signalStream:
			modEvStreams[sig.Module] <- EV_EXEC
		}
	}
}

// barUpdater updates the status-bar as module updates are submitted.
// To prevent very frequent module updates, the update submissions are batched
// and run in batches of at least `b.UpdateBatchMin` OR once every
// `b.UpdateMinTime` if there are any updates (which ever comes first).
//
// If `b.UpdateBatchMin` is 1, then the bar is updated whenever a module update
// is submitted.
func (b *StatusBar) barUpdater(
	outs []string,
	updateNotifyStream chan moduleUpdateNotification,
) {
	var batchedUpdates int // number of batched module updates
	for {
		select {
		case update := <-updateNotifyStream:
			outs[update.idx] = update.out
			batchedUpdates++
			// after receiving an update, if we have enough batched updates,
			// update the bar
			if batchedUpdates == b.UpdateBatchSize {
				batchedUpdates = 0
				b.collateAndUpdateBar(outs)
			}
		case <-time.After(b.UpdateMinTime):
			// after UpdateMinTime, if there are any batched updates, update the bar
			if batchedUpdates > 0 {
				batchedUpdates = 0
				b.collateAndUpdateBar(outs)
			}
		}
	}
}

// collateAndUpdateBar uses the given module outputs to form the status-bar
// string and calls xsetroot to update it.
func (b *StatusBar) collateAndUpdateBar(moduleOutputs []string) {
	// compose bar
	bar := " "
	for i, s := range moduleOutputs {
		bar += s
		if i != len(moduleOutputs)-1 {
			bar += b.Sep
		} else {
			bar += " "
		}
	}
	// update
	if err := xsetroot(bar); err != nil {
		log(ERROR, "%s", err)
	}
}
