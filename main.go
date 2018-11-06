package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/robfig/cron"
)

// Config File Structure, can be extend if needed, languages, update times, etc.
type Config struct {
	APIKey      string `json:"apikey"`
	CityIDTable []int  `json:"cityidtable"`
	CityIDx     int    `json:"cityidx"`
	Kindle      int    `json:"kindle"` // = 1 we are on kindle, do all kindle related jobs
	IconFont    string `json:"iconfont"`
	TxtFont     string `json:"txtfont"`
}

const numDays = 3   // number of forecast days to be displayed
const timeZones = 3 // number of time zones within one day e.g. morning, afternoon, evening

type displayTxtType struct {
	// Day       time.Time
	Description [numDays][timeZones]string
	Icon        [numDays][timeZones]string
	Temp        [numDays][timeZones]float64
	City        string
	TimeStamp   string
	Batt        string
}

func readConfig() (Config, error) {
	config := Config{APIKey: "",
		CityIDTable: []int{},
		CityIDx:     -1,
		Kindle:      0,
		IconFont:    "/usr/java/lib/fonts/kindleweathersr.ttf",
		TxtFont:     "/usr/java/lib/fonts/Helvetica_LT_65_Medium.ttf"}
	configFile, err := os.Open(conFile)
	defer configFile.Close()
	if err == nil {
		byteValue, _ := ioutil.ReadAll(configFile)
		err = json.Unmarshal(byteValue, &config)
		if err == nil {
			switch {
			case config.APIKey == "":
				err = errors.New("Undefined API Key")
			case config.CityIDx < 0:
				err = errors.New("Undifined City ID Index")
			case len(config.CityIDTable) == 0:
				err = errors.New("Undifined City ID Table, at least one ID needed")
			case (len(config.CityIDTable) - 1) < config.CityIDx:
				err = errors.New("City index bigger then the city numbers")
			}
		}
	}
	return config, err
}

// for the fisrt run add to all days and time zones icon "?" unknonw weather
func zeroDisplayTxt(displayTxt *displayTxtType) {

	for i := 0; i < numDays; i++ {
		for j := 0; j < timeZones; j++ {
			displayTxt.Icon[i][j] = "?"
		}
	}
}

func job(config *Config) {
	var err error

	mu.Lock()
	if config.Kindle == 1 {
		batLevel, _ := strconv.Atoi(checkBattery())
		displayTxt.Batt = ConverBatt(batLevel)
		err = CheckNetwork()
	} else {
		displayTxt.Batt = "?"
		err = nil
	}
	if err != nil {
		renderErrorDisp("!", err.Error())
	} else {
		w, err := getForecast5(config)
		if err != nil {
			log.Fatalln(err)
		}
		if w.City.ID != config.CityIDTable[config.CityIDx] { // open weather map did not return correct weather data, possible reason: network, server, etc. error?
			renderErrorDisp("!", "Weather data not available.")
		} else {
			ProcessWeatherData(config, &displayTxt, w)
		}
		// when we are on kidle go ahead with the display task
	}
	if config.Kindle == 1 {
		clearDisplay()
		showImage(picFile)
	}
	mu.Unlock()
}

var wg sync.WaitGroup
var mu sync.Mutex
var keyboard chan Kbd
var conFile = "/etc/yakwd.json"

const picFile = "/tmp/out.png"

// global display holds data from the past today's forecasts for morning, and afternoon, when they are not a forecast any loner
var displayTxt displayTxtType

func main() {

	if _, err := os.Stat(conFile); os.IsNotExist(err) {
		conFile = "yakwd.json"
	}
	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := os.Stat(config.IconFont); os.IsNotExist(err) {
		config.IconFont = "./fonts/kindleweathersr.ttf"
	}
	if _, err := os.Stat(config.TxtFont); os.IsNotExist(err) {
		config.TxtFont = "./fonts/Robotosr.ttf"
	}
	if config.Kindle != 1 {
		fmt.Println("Config: Not on Kindle!")
	}

	zeroDisplayTxt(&displayTxt)

	c := cron.New()
	c.AddFunc("@hourly", func() {
		job(&config)
	})
	c.Start()
	wg.Add(1)

	job(&config)

	if config.Kindle == 1 {
		keyboard = make(chan Kbd, 2)
		go KeyboardWorker()
		go MenuWorker(config)
		// wg.Add(1)
	}

	wg.Wait()

	os.Exit(0)
}
