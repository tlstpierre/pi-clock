module pi-clock

go 1.22.6

require (
	github.com/gorilla/mux v1.8.1
	github.com/sirupsen/logrus v1.9.3
	gopkg.in/yaml.v2 v2.4.0
	periph.io/x/conn/v3 v3.7.2
	periph.io/x/devices/v3 v3.7.2
	periph.io/x/host/v3 v3.8.2 // must be 3.8.2 or PWM won't work
)

require golang.org/x/sys v0.29.0 // indirect
