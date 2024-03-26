package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ht16k33"
	"time"
)

const (
	clockAddress = 0x70
)

var (
	i2cbus       i2c.BusCloser
	clockDevice  *ht16k33.Dev
	clockDisplay *NumDisplay
)

func i2cInit() {
	log.Debugf("Setting up i2c bus %s", ConfigData.I2cBus)
	var err error
	i2cbus, err = i2creg.Open(ConfigData.I2cBus)
	if err != nil {
		log.Fatalf("Problem opening i2c bus - %v", err)
	}
	clockDevice, err = ht16k33.NewI2C(i2cbus, clockAddress)
	if err != nil {
		log.Fatalf("Problem opening clock display - %v", err)
	}
	clockDevice.SetBrightness(15)
	clockDisplay = NewNumDisplay(clockDevice)
	clockDisplay.SetColon(false)

}

func i2cStop() {
	clockDevice.Halt()
	i2cbus.Close()
}

func updateLEDClock(currentTime time.Time) {
	hours, minutes, seconds := time.Now().Clock()
	clockDisplay.Write2Digits(0, uint8(hours), false)
	clockDisplay.Write2Digits(3, uint8(minutes), true)
	if seconds%2 == 1 {
		clockDisplay.SetColon(false)
	} else {
		clockDisplay.SetColon(true)
	}

	if hours < ConfigData.DisplayDimEnd || hours >= ConfigData.DisplayDimStart {
		clockDevice.SetBrightness(1)
	} else {
		clockDevice.SetBrightness(15)
	}

}

type NumDisplay struct {
	dev *ht16k33.Dev
}

func NewNumDisplay(device *ht16k33.Dev) *NumDisplay {
	return &NumDisplay{
		dev: device,
	}
}

func (d *NumDisplay) SetColon(state bool) {
	if state {
		d.dev.WriteColumn(2, 0)
	} else {
		d.dev.WriteColumn(2, 0xff)
	}
}

func (d *NumDisplay) WriteDigit(pos int, val uint8) error {
	if val > 16 {
		return fmt.Errorf("Value %d out-of-range for single digit", val)
	}
	err := d.dev.WriteColumn(pos, DigitMap[val])
	if err != nil {
		return err
	}
	return nil
}

func (d *NumDisplay) Write2Digits(pos int, val uint8, leadZero bool) error {
	val0 := val / 10
	var err error
	if val0 == 0 && !leadZero {
		log.Infof("Blanking position %d", pos)
		err = d.dev.WriteColumn(pos, 0)
	} else {
		err = d.dev.WriteColumn(pos, DigitMap[val0])
	}
	val1 := val % 10

	err = d.dev.WriteColumn(pos+1, DigitMap[val1])
	if err != nil {
		return err
	}
	return nil
}

var (
	DigitMap = map[uint8]uint16{
		0: 63, 1: 6, 2: 91, 3: 79, 4: 102, 5: 109, 6: 125, 7: 7, 8: 127, 9: 103,
		10: 119, 11: 124, 12: 88, 13: 94, 14: 121, 15: 113,
	}
)
