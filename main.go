package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/robfig/cron"
)

const picFile = "out.png"

// Config File Structure, can be extend if needed, languages, update times, etc.
type Config struct {
	APIKey string `json:"apikey"`
	CityID int    `json:"cityid"`
	Kindle int    `json:"kindle"` // = 1 we are on kindle, do all kindle related jobs
}

const numDays = 3   // number of forecast days to be displayed
const timeZones = 3 // number of time zones within one day e.g. morning, afternoon, evening

type displayTxtType struct {
	// Day       time.Time
	Description [timeZones]string
	Icon        [timeZones]string
	Temp        [timeZones]float64
}

func readConfig() (Config, error) {
	config := Config{APIKey: "", CityID: 0, Kindle: 0}
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err == nil {
		byteValue, _ := ioutil.ReadAll(configFile)
		err = json.Unmarshal(byteValue, &config)
		if err == nil {
			switch {
			case config.APIKey == "":
				err = errors.New("Undefined API Key")
			case config.CityID == 0:
				err = errors.New("Undifined City ID")
			}
		}
	}
	return config, err
}

// for the fisrt run add to all days and time zones icon "?" unknonw weather
func zeroDisplayTxt(displayTxt *[numDays]displayTxtType) {

	for i := 0; i < numDays; i++ {
		for j := 0; j < timeZones; j++ {
			displayTxt[i].Icon[j] = "?"
		}
	}
}

func job(config Config) {
	var err error

	if config.Kindle == 1 {
		err = CheckNetwork()
	} else {
		err = nil
	}
	if err != nil {
		renderErrorDisp("!", err.Error())
	} else {
		w, err := getForecast5(config)
		if err != nil {
			log.Fatalln(err)
		}
		if w.City.ID != config.CityID { // open weather map did not return correct weather data, possible reason: network, server, etc. error?
			renderErrorDisp("!", "Weather data not available.")
		} else {
			ProcessWeatherData(&displayTxt, w)
		}
		// when we are on kidle go ahead with the display task
		if config.Kindle == 1 {
			clearDisplay()
			showImage(picFile)
		}
	}
}

var wg sync.WaitGroup
var keyboard chan Kbd

// global display holds data from the past today's forecasts for morning, and afternoon, when they are not a forecast any loner
var displayTxt [numDays]displayTxtType

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if config.Kindle != 1 {
		fmt.Println("Config: Not on Kindle!")
	}

	zeroDisplayTxt(&displayTxt)

	c := cron.New()
	c.AddFunc("@hourly", func() {
		job(config)
	})
	c.Start()
	wg.Add(1)

	job(config)

	if config.Kindle == 1 {
		keyboard = make(chan Kbd, 2)
		go KeyboardWorker()
		go MenuWorker(config)
		// wg.Add(1)
	}

	wg.Wait()

	os.Exit(0)
}
