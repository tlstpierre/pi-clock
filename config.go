package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

type TimerAction string

const (
	FadeUp   TimerAction = "fadeup"
	FadeDown             = "fadedown"
	Bell                 = "bell"
)

type Config struct {
	TimeZone        string           `yaml:"timezone"`
	I2cBus          string           `yaml:"i2cbus"`
	PWMPin          string           `yaml:"pwmpin"`
	BrightLevel     int              `yaml:"brightlevel"`
	DimLevel        int              `yaml:"dimlevel"`
	DisplayDimStart int              `yaml:"displaydimstart"`
	DisplayDimEnd   int              `yaml:"displaydimend"`
	FadeUpTime      int              `yaml:"fadeuptime"`
	FadeDownTime    int              `yaml:"fadedowntime"`
	Timers          map[string]Timer `yaml:"timers"`
}

type Timer struct {
	Enabled   bool        `yaml:"enabled"`
	Sunday    bool        `yaml:"sunday"`
	Monday    bool        `yaml:"monday"`
	Tuesday   bool        `yaml:"tuesday"`
	Wednesday bool        `yaml:"wednesday"`
	Thursday  bool        `yaml:"thursday"`
	Friday    bool        `yaml:"friday"`
	Saturday  bool        `yaml:"saturday"`
	Hour      int         `yaml:"hour"`
	Minute    int         `yaml:"minute"`
	Action    TimerAction `yaml:"function"`
}

// Initialize a config object with default values
func (c *Config) Initialize() {
	*c = Config{
		TimeZone:        "America/Toronto",
		I2cBus:          "",
		PWMPin:          "",
		BrightLevel:     15,
		DimLevel:        1,
		DisplayDimStart: 21,
		DisplayDimEnd:   6,
		FadeUpTime:      30,
		FadeDownTime:    1,
		Timers: map[string]Timer{
			"Weekday": Timer{
				Enabled:   true,
				Hour:      07,
				Minute:    0,
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Action:    FadeUp,
			},
			"WeekdayOff": Timer{
				Enabled:   true,
				Hour:      8,
				Minute:    0,
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Action:    FadeDown,
			},

			"Saturday": Timer{
				Hour:     9,
				Minute:   15,
				Saturday: true,
				Enabled:  true,
				Action:   FadeUp,
			},
			"SaturdayOff": Timer{
				Hour:     10,
				Minute:   0,
				Saturday: true,
				Enabled:  true,
				Action:   FadeDown,
			},

			"Sunday": Timer{
				Hour:    7,
				Minute:  0,
				Sunday:  true,
				Enabled: true,
				Action:  FadeUp,
			},
			"SundayOff": Timer{
				Hour:    8,
				Minute:  0,
				Sunday:  true,
				Enabled: true,
				Action:  FadeDown,
			},
		},
	}
}

// Load any YAML values into the config object from a file.
func (c *Config) LoadFile(filename string) error {
	// Open config file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&c); err != nil {
		return err
	}
	return nil
}

func (c *Config) SaveFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	e := yaml.NewEncoder(file)
	// Start YAML encoding to file
	if err := e.Encode(c); err != nil {
		return err
	}
	return nil

}
