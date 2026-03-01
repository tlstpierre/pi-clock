package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"periph.io/x/host/v3"

	//	"periph.io/x/host/v3/rpi"

	"syscall"
	"time"
)

const (
	tickerPrecision = 3500000
)

var (
	LogLevel        = flag.String("loglevel", "info", "Log Level")
	ConfigFile      = flag.String("configfile", "config.yaml", "A yaml config file to use")
	ConfigData      = &Config{}
	TZLocation      *time.Location
	secondTicker    *time.Ticker
	timeError       int
	timeCalibrating bool
	Dimmer1         *Dimmer
)

func main() {
	//	Set the default debug lever
	// Parse the runtime flags
	flag.Parse()
	lvl, _ := log.ParseLevel(*LogLevel)
	log.SetLevel(lvl)

	ConfigData.Initialize()
	err := ConfigData.LoadFile(*ConfigFile)
	if err != nil {
		log.Errorf("Could not open config file %v", err)
	}
	log.Infof("Config contents are %+v", ConfigData)
	TZLocation, _ = time.LoadLocation(ConfigData.TimeZone)

	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	signal.Notify(sigs)
	setTicker()

	// Set up our peripherals
	if hostState, err := host.Init(); err != nil {
		log.Fatalf("Problem setting up periph host - %v", err)
	} else {
		log.Debugf("Periph host state is %+v", hostState)
	}

	Dimmer1, err = NewDimmer(ConfigData.PWMPin)
	if err != nil {
		log.Fatalf("Problem setting up dimmer - %v", err)
	}
	defer Dimmer1.Stop()
	i2cInit()

	log.Info("Set up dimmer")
	Dimmer1.Fade(FullBright, 2*time.Second)
	time.Sleep(4 * time.Second)
	Dimmer1.Fade(0, 2*time.Second)

	StartServer()

	for {
		select {
		case s := <-sigs:
			switch s {
			case syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT:
				AppCleanup()
				os.Exit(1)
			case syscall.SIGHUP:

				//			case syscall.SIGINFO:
			}
		case currentTime := <-secondTicker.C:
			calTicker(currentTime)
			//			log.Debugf("Time is %v with error %d", currentTime, timeError)
			updateLEDClock(currentTime)
			if currentTime.Second() == 0 {
				action := checkTimer(currentTime)
				if action != nil {
					switch *action {
					case FadeUp:
						Dimmer1.Fade(FullBright, time.Duration(ConfigData.FadeUpTime)*time.Minute)
					case FadeDown:
						Dimmer1.Fade(0, time.Duration(ConfigData.FadeDownTime)*time.Minute)
					}
				}
			}
		}
	}
}

func AppCleanup() {
	i2cStop()
	StopServer()
}

func setTicker() {
	nextSecond := time.Now().Truncate(time.Second)
	timeWait := time.Until(nextSecond)
	time.Sleep(timeWait)
	secondTicker = time.NewTicker(time.Second)
	log.Infof("Started ticker at %v", time.Now().Nanosecond())

}

func calTicker(currentTime time.Time) {
	nsError := currentTime.Nanosecond()
	if nsError > 500000000 {
		timeError = 1000000000 - nsError
	} else {
		timeError = nsError
	}
	if timeError > tickerPrecision*2 || timeError < -tickerPrecision/2 {
		log.Info("Calibrating ticker")
		newInterval := time.Second - time.Duration(timeError)/2
		secondTicker.Reset(newInterval)
		timeCalibrating = true
	} else if timeCalibrating && (timeError < tickerPrecision || timeError < -tickerPrecision) {
		secondTicker.Reset(time.Second)
		timeCalibrating = false
		log.Info("Ticker calibrated")
	}

}
