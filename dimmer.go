package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"time"
)

const (
	lampIO     = "PWM1"
	FullBright = 16777216
)

type Dimmer struct {
	pin       gpio.PinIO
	intensity int32
	target    int32
	fadetime  time.Duration
	freq      physic.Frequency
	newValue  chan struct{}
	stopChan  chan struct{}
}

func NewDimmer(pin string) (dimmer *Dimmer, err error) {
	dimmer = &Dimmer{
		pin:      gpioreg.ByName(pin),
		fadetime: time.Second,
		freq:     100 * physic.Hertz,
		newValue: make(chan struct{}),
		stopChan: make(chan struct{}),
	}
	if dimmer.pin == nil {
		err = fmt.Errorf("Could not register pin %s", pin)
		return
	}
	log.Debugf("Pin %s is set to %s", pin, dimmer.pin.Function())

	// Set the initial pin state
	if err = dimmer.pin.Out(gpio.Low); err != nil {
		err = fmt.Errorf("Problem setting pin %s to output low - %v", pin, err)
		return
	}

	// Run a dimmer goroutine
	go dimmerRoutine(dimmer)
	return
}

func (d *Dimmer) Intensity(intensity int32) error {
	d.intensity = intensity
	d.target = intensity
	d.newValue <- struct{}{}
	return nil
}

func (d *Dimmer) Fade(intensity int32, rate time.Duration) error {
	d.target = intensity
	d.fadetime = rate
	d.newValue <- struct{}{}
	return nil
}

// Clean shutdown
func (d *Dimmer) Stop() {
	d.stopChan <- struct{}{}
	d.pin.Halt()
}

func dimmerRoutine(dimmer *Dimmer) {
	log.Infof("Starting dimmer routine on pin %s", dimmer.pin)
	//		var isfading bool
	fadeTicker := time.NewTicker(time.Hour)
	//var fadeTicker *time.Ticker
	var stepSize int32

	for {
		select {
		case <-dimmer.newValue:
			log.Debug("Got change in values for dimmer")
			dimmer.setOutput()
			if dimmer.target != dimmer.intensity {
				log.Debugf("Dimmer fade time is %v", dimmer.fadetime)
				stepCount := dimmer.intensity - dimmer.target
				if stepCount < 0 {
					stepCount = -stepCount
					log.Debugf("Fading up %d steps", stepCount)
				} else {
					log.Debugf("Fading down %d steps", stepCount)
				}
				interval := dimmer.fadetime / time.Duration(stepCount)
				log.Debugf("Calculated interval is %v", interval)
				if interval < (50 * time.Millisecond) {
					steps := dimmer.fadetime.Milliseconds() / 100
					log.Debugf("Dimmer interval is too short - scaling to %d steps at 100ms", steps)
					interval = 100 * time.Millisecond
					log.Debugf("Setting step size to %d / %d", steps, stepCount)
					stepSize = stepCount / int32(steps)
				} else {
					stepSize = 1
				}
				log.Debugf("step size is %d", stepSize)
				log.Debugf("Starting fade from %d to %d for %d steps at interval %v", dimmer.intensity, dimmer.target, stepCount, interval)
				fadeTicker = time.NewTicker(interval)
				//					isfading = true
			}
		case <-fadeTicker.C:
			log.Debugf("Fading - %d to %d step size %d", dimmer.intensity, dimmer.target, stepSize)
			if dimmer.target > dimmer.intensity {
				dimmer.intensity += stepSize
			} else if dimmer.target < dimmer.intensity {
				dimmer.intensity -= stepSize
			}
			dimmer.setOutput()
			gap := dimmer.target - dimmer.intensity
			if gap < 0 {
				gap = -gap
			}
			if gap == 0 || gap < stepSize {
				if gap > 0 {
					dimmer.intensity = dimmer.target
					dimmer.setOutput()
				}
				log.Debugf("Reached target of %d", dimmer.target)
				// fadeTicker.Reset(time.Hour)
				fadeTicker.Stop()
				//			continue
			}

		case <-dimmer.stopChan:
			log.Debug("Stopping dimmer")
			return
		}
	}
}

func (dimmer *Dimmer) setOutput() {
	log.Infof("Setting dimmer to intensity %d out of %d", dimmer.intensity, gpio.DutyMax)
	if dimmer.intensity == 0 {
		dimmer.pin.Out(gpio.Low)
	} else if gpio.Duty(dimmer.intensity) >= gpio.DutyMax {
		dimmer.pin.Out(gpio.High)
	} else {
		duty := gpio.Duty(dimCurve(dimmer.intensity))
		if duty.Valid() {
			err := dimmer.pin.PWM(duty, dimmer.freq)
			if err != nil {
				log.Errorf("Could not set pin %s to intensity %d - %v", dimmer.pin, dimmer.intensity, err)
			}
		} else {
			log.Errorf("Duty %v is invalid", duty)
		}
		log.Debugf("Duty cycle is %s", duty)
	}

}

func dimCurve(input int32) int32 {
	output := math.Pow(float64(input), 1.5) / 4096
	return int32(output)
}
