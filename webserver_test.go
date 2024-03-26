package main

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestWebServer(t *testing.T) {
	ConfigData.Initialize()
	/*
		err := ConfigData.LoadFile(*ConfigFile)
		if err != nil {
			log.Errorf("Could not open config file %v", err)
		}
	*/
	log.Infof("Config contents are %+v", ConfigData)
	TZLocation, _ = time.LoadLocation(ConfigData.TimeZone)
	StartServer()
	time.Sleep(30 * time.Second)
	StopServer()
}
