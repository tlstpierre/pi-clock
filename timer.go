package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

var ()

func checkTimer(currentTime time.Time) *TimerAction {
	for timerName, timerData := range ConfigData.Timers {
		log.Debugf("Checking timer %s", timerName)
		var alarmTime time.Time
		switch timerData.Action {
		case FadeUp:
			alarmTime = currentTime.Add(time.Duration(ConfigData.FadeUpTime) * time.Minute)
		case FadeDown:
			alarmTime = currentTime.Add(time.Duration(ConfigData.FadeDownTime) * time.Minute)
		default:
			alarmTime = currentTime
		}
		log.Infof("Checking for alarm at time %v so we can start fading", alarmTime)
		// Ignore if disabled
		if !timerData.Enabled {
			continue
		}
		// Check the day of the week
		switch alarmTime.Weekday() {
		case time.Sunday:
			if !timerData.Sunday {
				continue
			}
		case time.Monday:
			if !timerData.Monday {
				continue
			}
		case time.Tuesday:
			if !timerData.Tuesday {
				continue
			}
		case time.Wednesday:
			if !timerData.Wednesday {
				continue
			}
		case time.Thursday:
			if !timerData.Thursday {
				continue
			}
		case time.Friday:
			if !timerData.Friday {
				continue
			}
		case time.Saturday:
			if !timerData.Saturday {
				continue
			}
		}
		log.Debugf("Timer %s matches current day", timerName)
		// Check the hour
		if alarmTime.Hour() != timerData.Hour {
			continue
		}
		log.Debugf("Timer %s matches current hour", timerName)

		// Check the minute
		if alarmTime.Minute() != timerData.Minute {
			continue
		}
		// Fire our timer
		log.Infof("Timer %s matches", timerName)
		return &timerData.Action
	}
	return nil
}
