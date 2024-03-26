package main

import (
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	//	"html/template"
	"net/http"
	"strconv"
	"time"
)

var (
	Server *http.Server
)

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", HandleMainForm).Methods("GET")
	r.HandleFunc("/", HandleMainFormPost).Methods("POST")

	Server = &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	_, err := MainForm.Parse(webform)
	if err != nil {
		log.Fatalf("Could not parse webform - %v", err)
	}

	go func() {
		log.Fatal(Server.ListenAndServe())
	}()
}

func StopServer() {
	log.Info("Stopping webserver")
	Server.Shutdown(context.TODO())
}

func HandleMainForm(w http.ResponseWriter, r *http.Request) {
	formData := FormData{
		CurrentTime: time.Now(),
		Hours:       hours,
		Minutes:     minutes,
		Actions:     actions,
		Timers:      ConfigData.Timers,
	}
	log.Info("Executing main form")
	err := MainForm.Execute(w, formData)
	if err != nil {
		log.Errorf("Problem executing main form - %v", err)
	}
}

func HandleMainFormPost(w http.ResponseWriter, r *http.Request) {
	log.Infof("Submit action is %s", r.PostFormValue("SubmitAction"))
	var fileChanged bool
	switch r.PostFormValue("SubmitAction") {
	case "Set":
		value := ParseInt(r.PostFormValue("level"))
		log.Infof("Manually setting level to %d", value)
		Dimmer1.Fade(int32(value)*16777216/100, 3*time.Second)
	case "Update", "Create":
		timerName := r.PostFormValue("timer")
		if timerName == "" {
			log.Errorf("Form posted with no timer name")
			return
		}
		if r.PostFormValue("SubmitAction") == "Update" {
			_, timerExists := ConfigData.Timers[timerName]
			if !timerExists {
				log.Errorf("Timer %s does not exist already", timerName)
				break
			}
		}
		newTimer := Timer{
			Enabled:   ParseBool(r.PostFormValue("enabled")),
			Hour:      ParseInt(r.PostFormValue("hour")),
			Minute:    ParseInt(r.PostFormValue("minute")),
			Sunday:    ParseBool(r.PostFormValue("Sunday")),
			Monday:    ParseBool(r.PostFormValue("Monday")),
			Tuesday:   ParseBool(r.PostFormValue("Tuesday")),
			Wednesday: ParseBool(r.PostFormValue("Wednesday")),
			Thursday:  ParseBool(r.PostFormValue("Thursday")),
			Friday:    ParseBool(r.PostFormValue("Friday")),
			Saturday:  ParseBool(r.PostFormValue("Saturday")),
			Action:    TimerAction(r.PostFormValue("Action")),
		}
		log.Infof("New timer is %+v", newTimer)
		ConfigData.Timers[timerName] = newTimer
		fileChanged = true
	case "Delete":
		timerName := r.PostFormValue("timer")
		if timerName == "" {
			log.Errorf("Form posted with no timer name")
			return
		}
		delete(ConfigData.Timers, timerName)
		fileChanged = true
	}

	if fileChanged {
		log.Infof("Saving new config file %s", *ConfigFile)
		err := ConfigData.SaveFile(*ConfigFile)
		if err != nil {
			log.Errorf("Could not save config - %v", err)
		}
	}
	http.Redirect(w, r, r.URL.Path, 301)
}

func ParseBool(input string) bool {
	if input == "on" {
		return true
	}
	return false
}

func ParseInt(input string) int {
	value, _ := strconv.ParseInt(input, 10, 64)
	return int(value)
}
