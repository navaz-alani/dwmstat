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

func (b *StatusBar) Exec() {
	signalStream := make(chan Signal, 5)
	go listenAndNotify(signalStream)
	b.Udpate(signalStream)
}

// Event is a kind of signal sent to a worker.
type Event string

const (
	EV_STOP = "STOP" // skip current update internal and stop
	EV_EXEC = "EXEC" // don't wait for interval to finish, run now
)

type ModuleUpdateCallback func(out string)

// moduleUpdater is the function which periodically executes a module and
// submits its output.
// It can also receive and process events through the evStream, so if a module
// is not to be periodically updated, it can be updated through a signal
// instead.
func moduleUpdater(mod Module, evStream chan Event, cb ModuleUpdateCallback) {
	if mod.UpdateInterval() > 0 {
		select {
		case ev := <-evStream:
			switch ev {
			case EV_EXEC: // run without waiting for interval to complete
				log(INFO, "EV_EXEC - updating module '%s'", mod.Name())
				cb(mod.Exec())
			case EV_STOP:
				return
			}
		case <-time.After(mod.UpdateInterval()):
			log(INFO, "UPDATE_INTERVAL - updating module '%s'", mod.Name())
			cb(mod.Exec())
		}
	} else {
		ev := <-evStream
		switch ev {
		case EV_EXEC:
			log(INFO, "EV_EXEC - updating module '%s'", mod.Name())
			cb(mod.Exec())
		case EV_STOP:
			return
		}
	}
	go moduleUpdater(mod, evStream, cb) // respawn
}

type moduleUpdateNotification struct {
	idx int
	out string
}

// Update updates the status bar.
func (b *StatusBar) Udpate(signalStream chan Signal) {
	// mapping of module name to Module
	modules := make(map[string]Module)
	// channels through which modules receive events
	moduleEventStreams := make(map[string]chan Event)
	// mapping between a signal and names of modules to update
	signals := make(map[int][]string)

	// collect all modules
	for _, m := range b.Modules {
		mod := m.Name()
		if _, ok := modules[mod]; ok {
			log(WARN, "duplicate module name '%s'", mod)
		}
		modules[mod] = m
		moduleEventStreams[mod] = make(chan Event)
		signals[m.UpdateSig()] = append(signals[m.UpdateSig()], mod)
	}

	outs := make([]string, len(b.Modules)) // outputs produced by modules

	// initially run all modules
	for i, m := range b.Modules {
		outs[i] = m.Exec()
	}

	var updateNotifyStream = make(chan moduleUpdateNotification)

	// dispatch routines for each module
	for i, mod := range b.Modules {
		i := i
		go moduleUpdater(
			mod,
			moduleEventStreams[mod.Name()],
			func(o string) { updateNotifyStream <- moduleUpdateNotification{i, o} },
		)
	}

	// routine that collects module updates and updates the status bar
	go b.barUpdater(outs, updateNotifyStream)

	for {
		select {
		case sig := <-signalStream:
			moduleEventStreams[sig.Module] <- EV_EXEC
		}
	}
}

// barUpdater updates the status-bar as module updated are submitted.
// To prevent very frequent module updates, the update submissions are batched
// and run in batches of at least `b.UpdateBatchMin` OR once every
// `b.UpdateMinTime` if there are any updates (which ever comes first).
func (b *StatusBar) barUpdater(
	outs []string,
	updateNotifyStream chan moduleUpdateNotification,
) {
	var batchedUpdates int
	for {
		select {
		case update := <-updateNotifyStream:
			outs[update.idx] = update.out
			batchedUpdates++
			if batchedUpdates == b.UpdateBatchSize {
				batchedUpdates = 0
				b.collateAndUpdateBar(outs)
			}
		case <-time.After(b.UpdateMinTime):
			if batchedUpdates > 0 {
				batchedUpdates = 0
				b.collateAndUpdateBar(outs)
			}
		}
	}
}

func (b *StatusBar) collateAndUpdateBar(moduleOutputs []string) {
	bar := " "
	for i, s := range moduleOutputs {
		bar += s
		if i != len(moduleOutputs)-1 {
			bar += b.Sep
		} else {
			bar += " "
		}
	}

	if err := xsetroot(bar); err != nil {
		log(ERROR, "%s", err)
	}
}
