package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.com/timstpierre/go-embedded/pkg/lcd1602"
	"gitlab.com/timstpierre/go-embedded/pkg/seesaw"
	"sync"
	"time"
)

var (
	BacklightTime = 1 * time.Minute
)

type DisplayManager struct {
	sync.RWMutex
	backlightState   bool
	daytime          bool
	stale            bool
	currentIntensity int32
	displayMode      DisplayMode
	lcd              *lcd1602.Dev
	encoder          *seesaw.Dev
	dimmer           *Dimmer
	screenTimer      *time.Timer
	inputPoll        *time.Ticker
	refreshTimer     *time.Ticker
	context          context.Context
	cancelFunc       context.CancelFunc
}

type DisplayMode uint8

const (
	Status DisplayMode = iota
	TimerList
	EditTimer
	Weather
	Settings
)

func NewDisplayManager(lcdDevice *lcd1602.Dev, encoderDevice *seesaw.Dev, dimmer *Dimmer) *DisplayManager {
	mgr := &DisplayManager{
		lcd:          lcdDevice,
		encoder:      encoderDevice,
		dimmer:       dimmer,
		screenTimer:  time.NewTimer(BacklightTime),
		inputPoll:    time.NewTicker(100 * time.Millisecond),
		refreshTimer: time.NewTicker(5 * time.Second),
		stale:        true,
	}
	mgr.context, mgr.cancelFunc = context.WithCancel(context.Background())
	go mgr.run()

	return mgr
}

func (d *DisplayManager) Stop() {
	d.cancelFunc()
}

func (d *DisplayManager) run() {
	log.Info("Starting display manager")
	d.lcd.SetBacklight(true)
	d.lcd.Home()
	d.lcd.Clear()
	d.lcd.CursorMode(false, false)

	for {
		select {
		case <-d.context.Done():
			// Exit gracefully
			busLock.Lock()
			d.lcd.Clear()
			d.lcd.SetBacklight(false)
			busLock.Unlock()
			_ = d.screenTimer.Stop()
			d.inputPoll.Stop()
			return

		case <-d.refreshTimer.C:
			d.refresh()

		case <-d.screenTimer.C:
			// Turn off the screen
			busLock.Lock()
			d.lcd.SetBacklight(false)
			busLock.Unlock()
			d.backlightState = false

		case <-d.inputPoll.C:
			// Read the encoder
		}
	}
}

func (d *DisplayManager) Day() {
	d.RLock()
	daytime := d.daytime
	d.RUnlock()
	if daytime {
		return
	}
	d.Lock()
	defer d.Unlock()

	d.screenTimer.Stop()
	d.daytime = true
	d.lcd.SetBacklight(true)
}

func (d *DisplayManager) Night() {
	d.RLock()
	daytime := d.daytime
	d.RUnlock()
	if !daytime {
		return
	}
	d.Lock()
	defer d.Unlock()

	d.screenTimer.Reset(BacklightTime)
	d.daytime = false
	d.lcd.SetBacklight(true)
}

func (d *DisplayManager) SetDimmer(dim *Dimmer) {
	d.dimmer = dim
}

func (d *DisplayManager) refresh() {
	switch d.displayMode {
	case Status:
		if d.dimmer != nil {
			intensity, target := d.dimmer.PercentOutput()
			_, _, hour, minute := nextTimer(FadeUp)
			busLock.Lock()
			time.Sleep(100 * time.Millisecond)
			d.lcd.WriteLineUTFString(1, fmt.Sprintf("Lvl: %d%% / %d%%", intensity, target))
			d.lcd.WriteLineUTFString(2, fmt.Sprintf("Nxt: %02d:%02d", hour, minute))
			time.Sleep(100 * time.Millisecond)
			busLock.Unlock()
		} else {
			busLock.Lock()
			d.lcd.WriteUTFString("Dimmer not ready")
			busLock.Unlock()
		}
	}
	d.stale = false

}
