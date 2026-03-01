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

func nextTimer(action TimerAction) (name, day string, hour, minute int) {
	currentHour := time.Now().Hour()
	currentMinute := time.Now().Minute()
	currentDay := time.Now().Weekday()
	// next timer today
	for timerName, timerData := range ConfigData.Timers {
		log.Debugf("Checking to see if %s is next today", timerName)
		if !timerData.Enabled {
			continue
		}
		if timerData.Action != action {
			continue
		}
		switch currentDay {
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
		if hour == 0 && timerData.Hour >= currentHour && timerData.Hour <= hour && timerData.Minute >= currentMinute && timerData.Minute < minute {
			hour = timerData.Hour
			minute = timerData.Minute
			name = timerName
		} else if timerData.Hour >= currentHour && timerData.Hour <= hour {
			if timerData.Minute >= currentMinute && timerData.Minute < minute {
				log.Debugf("Timer %s is sooner than %s", timerName, name)
				hour = timerData.Hour
				minute = timerData.Minute
				name = timerName
			}
		}
	}

	// _Found a match today
	if name != "" {
		return
	} else {
		log.Debug("Nothing matching today, try tomorrow")
	}

	for timerName, timerData := range ConfigData.Timers {
		log.Debugf("Checking to see if %s is next tomorrow", timerName)
		if !timerData.Enabled {
			continue
		}
		if timerData.Action != action {
			continue
		}
		switch currentDay {
		case time.Saturday:
			if !timerData.Sunday {
				continue
			}
		case time.Sunday:
			if !timerData.Monday {
				continue
			}
		case time.Monday:
			if !timerData.Tuesday {
				continue
			}
		case time.Tuesday:
			if !timerData.Wednesday {
				continue
			}
		case time.Wednesday:
			if !timerData.Thursday {
				continue
			}
		case time.Thursday:
			if !timerData.Friday {
				continue
			}
		case time.Friday:
			if !timerData.Saturday {
				continue
			}
		}
		if hour == 0 {
			log.Debugf("No timer found yet - putting %s down as placeholder", timerName)
			hour = timerData.Hour
			minute = timerData.Minute
			name = timerName
		} else if timerData.Hour <= hour {
			if timerData.Minute < minute {
				log.Debugf("Timer %s is sooner than %s", timerName, name)
				hour = timerData.Hour
				minute = timerData.Minute
				name = timerName
			}
		}
	}

	return
}
