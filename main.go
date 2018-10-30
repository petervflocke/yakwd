package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/fogleman/gg"
	"github.com/robfig/cron"
)

const picFile = "out.png"

// Config File Structure, can be extend if needed, languages, update times, etc.
type Config struct {
	APIKey string `json:"apikey"`
	CityID int    `json:"cityid"`
	Kindle int    `json:"kindle"` // = 1 we are on kindle, do all kindle related jobs
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

const numDays = 3   // number of forecast days to be displayed
const timeZones = 3 // number of time zones within one day e.g. morning, afternoon, evening

type displayTxtType struct {
	// Day       time.Time
	Description [timeZones]string
	Icon        [timeZones]string
	Temp        [timeZones]float64
}

func renderWeatherDisp(displayTxt *[numDays]displayTxtType) {
	const my = 75       // margin from top
	const dy = 25 + 200 // high of a column
	const ix = 70       // x of a midle of a weather icon in the firts column
	const dx = 200      // width of a column
	const wy = 70       // y of a midle of a weather icon in a first column
	const tx = 160      // x of a first temperature text
	const ty = 115      // y of a first temperature text
	const hy = -20      // move up header text
	const iy = 190      // y of the the detailed weather information

	// for i := 0; i < numDays; i++ {
	// 	for j := 0; j < timeZones; j++ {
	// 		fmt.Print(displayTxt[i].Temp[j], ", ")
	// 	}
	// 	fmt.Println("")
	// }

	dc := gg.NewContext(600, 800)
	//img := image.NewGray(image.Rect(0, 0, 600, 800))
	//dc := gg.NewContextForImage(img) // does not work, img will be anyway converted to full RGB and requires closing convertion into gray scale
	ClearPic(dc)

	// create a table 5x3 (header, forecast 3x3, footer)
	dc.SetLineWidth(1)
	for i := 0; i < 4; i++ {
		dc.DrawLine(0, my+float64(i*dy), 600, my+float64(i*dy))
	}
	dc.DrawLine(1*dx, 0, 1*dx, 800)
	dc.DrawLine(2*dx, 0, 2*dx, 800)

	for i := 0; i < numDays; i++ {
		// Print Temperature
		if err := dc.LoadFontFace("./fonts/Robotosr.ttf", 70); err != nil {
			panic(err)
		}
		for j := 0; j < timeZones; j++ {
			if displayTxt[i].Icon[j] != "?" {
				dc.DrawStringAnchored(fmt.Sprintf("%-0.0f", displayTxt[i].Temp[j]), tx+float64(j)*dx, my+ty+float64(i)*dy, 1, 1)
			}
		}

		// add smaller °C at the end of each temp.
		if err := dc.LoadFontFace("./fonts/Robotosr.ttf", 30); err != nil {
			panic(err)
		}
		for j := 0; j < timeZones; j++ {
			if displayTxt[i].Icon[j] != "?" {
				dc.DrawStringAnchored("°C", tx+float64(j)*dx, my+ty+float64(i)*dy, 0, 1)
			}
		}

		// add detail (small text) weather description
		if err := dc.LoadFontFace("./fonts/Robotosr.ttf", 16); err != nil {
			panic(err)
		}
		for j := 0; j < timeZones; j++ {
			if displayTxt[i].Icon[j] != "?" {
				dc.DrawStringWrapped(displayTxt[i].Description[j], dx/2+float64(j)*dx, my+iy+float64(i)*dy, 0.5, 0.5, 180, 1.5, gg.AlignCenter)
			}
		}

		// print weather icon
		if err := dc.LoadFontFace("./fonts/kindleweathersr.ttf", 100); err != nil {
			panic(err)
		}

		for j := 0; j < timeZones; j++ {
			dc.DrawStringAnchored(displayTxt[i].Icon[j], ix+float64(j)*dy, my+wy+float64(i)*dy, 0.5, 0.5)
		}

	}
	// other static text
	if err := dc.LoadFontFace("./fonts/Robotosr.ttf", 30); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Morning", dx/2+0*dx, my+hy, 0.5, 0)
	dc.DrawStringAnchored("Afternoon", dx/2+1*dx, my+hy, 0.5, 0)
	dc.DrawStringAnchored("Evening", dx/2+2*dx, my+hy, 0.5, 0)

	dc.Stroke()
	// dc.SavePNG("tmp.png") // not a gray scale picture - we do not need it
	SaveGrayPic(dc.Image(), picFile)
}

func renderErrorDisp(message string) {
	const my = 75 // margin from top

	dc := gg.NewContext(600, 800)
	ClearPic(dc)

	// print ! icon
	if err := dc.LoadFontFace("./fonts/kindleweathersr.ttf", 200); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("!", 300, my+200, 0.5, 0.5)

	if err := dc.LoadFontFace("./fonts/Robotosr.ttf", 40); err != nil {
		panic(err)
	}
	dc.DrawStringWrapped(message, 300, 450, 0.5, 0.5, 550, 1.5, gg.AlignCenter)
	dc.Stroke()

	SaveGrayPic(dc.Image(), picFile)
}

// for the fisrt run add to all days and time zones icon "?" unknonw weather
func zeroDisplayTxt(displayTxt *[numDays]displayTxtType) {

	for i := 0; i < numDays; i++ {
		for j := 0; j < timeZones; j++ {
			displayTxt[i].Icon[j] = "?"
		}
	}
}

func processWeatherData(displayTxt *[numDays]displayTxtType, w *owm.Forecast5WeatherData) {

	var weatherFontMapping = map[int]string{
		// openweather Main weather code to icon conversion
		200: "s", 201: "s", 202: "s", 210: "p", 211: "s", 212: "s", 221: "s", 230: "s", 231: "s", 232: "s",
		300: "k", 301: "k", 302: "k", 310: "k", 311: "k", 312: "k", 313: "k", 314: "k", 321: "k",
		500: "c", 501: "c", 502: "c", 503: "b", 504: "b", 511: "i", 520: "c", 521: "c", 522: "c", 531: "c",
		600: "r", 601: "r", 602: "r", 611: "f", 612: "r", 615: "u", 616: "u", 620: "r", 621: "r", 622: "r",
		701: "t", 711: "o", 721: "h", 731: "d", 741: "h", 751: "d", 761: "d", 762: "d", 771: "v", 781: "w",
		800: "n", 801: "j", 802: "j", 803: "a", 804: "x"}

	currentTime := time.Now().UTC()
	localTime := currentTime
	location, err := time.LoadLocation("Local")
	if err == nil {
		localTime = localTime.In(location)
	}

	// fmt.Println("UTC: ", currentTime)
	// fmt.Println("Loc: ", localTime)
	currentTime = localTime

	// time zones ranges are defined based on 3 hours blocks as the data are comming from open weather maps
	// Starting from 00:00, 03:00, 06:00, 09:00, 12, 15, 18, 21, 24=00:00
	// Current Check Points: 6=morning, 12=afternoon, 18-evening
	// Check points are compare with the open weather map time in UTC format !!!
	d1TimeZone1 := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 6, 0, 0, 0, time.UTC)
	d1TimeZone2 := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 12, 0, 0, 0, time.UTC)
	d1TimeZone3 := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 18, 0, 0, 0, time.UTC)
	d2TimeZone1 := d1TimeZone1.AddDate(0, 0, 1) // tomorrow
	d2TimeZone2 := d1TimeZone2.AddDate(0, 0, 1)
	d2TimeZone3 := d1TimeZone3.AddDate(0, 0, 1)
	d3TimeZone1 := d1TimeZone1.AddDate(0, 0, 2) // the day after tomorrow
	d3TimeZone2 := d1TimeZone2.AddDate(0, 0, 2)
	d3TimeZone3 := d1TimeZone3.AddDate(0, 0, 2)

	// Convert all to unix int format to compare it with dt time from open weather maps
	d1T1U := d1TimeZone1.Unix() // day 1 Time zone 1 format Unix
	d1T2U := d1TimeZone2.Unix() // day 2 Time zone 2 format Unix
	d1T3U := d1TimeZone3.Unix() // ...
	d2T1U := d2TimeZone1.Unix()
	d2T2U := d2TimeZone2.Unix()
	d2T3U := d2TimeZone3.Unix()
	d3T1U := d3TimeZone1.Unix()
	d3T2U := d3TimeZone2.Unix()
	d3T3U := d3TimeZone3.Unix()

	var key owm.Forecast5WeatherList
	for i := len(w.List) - 1; i >= 0; i-- {

		forecastDay := -1
		timeZone := -1
		key = w.List[i]

		switch tt := key.Dt; {
		case (tt >= d1T1U) && (tt < d1T2U):
			forecastDay = 0
			timeZone = 0
		case (tt >= d1T2U) && (tt < d1T3U):
			forecastDay = 0
			timeZone = 1
		case (tt >= d1T3U) && (tt < d2T1U):
			forecastDay = 0
			timeZone = 2
		case (tt >= d2T1U) && (tt < d2T2U):
			forecastDay = 1
			timeZone = 0
		case (tt >= d2T2U) && (tt < d2T3U):
			forecastDay = 1
			timeZone = 1
		case (tt >= d2T3U) && (tt < d3T1U):
			forecastDay = 1
			timeZone = 2
		case (tt >= d3T1U) && (tt < d3T2U):
			forecastDay = 2
			timeZone = 0
		case (tt >= d3T2U) && (tt < d3T3U):
			forecastDay = 2
			timeZone = 1
		case (tt >= d3T3U):
			forecastDay = 2
			timeZone = 2
		}
		if (forecastDay != -1) && (timeZone != -1) {
			displayTxt[forecastDay].Temp[timeZone] = key.Main.Temp
			displayTxt[forecastDay].Description[timeZone] = key.Weather[0].Description
			displayTxt[forecastDay].Icon[timeZone] = weatherFontMapping[key.Weather[0].ID]
			if displayTxt[forecastDay].Icon[timeZone] == "" {
				displayTxt[forecastDay].Icon[timeZone] = "?"
			}
			//customLogf("Day: %d, Time: %d  ", forecastDay, timeZone)
			//fmt.Println("dt:", time.Unix(key.Dt, 0).UTC(), " dt: ", key.Dt)
		}
	}
	renderWeatherDisp(displayTxt)
}

func getForecast5(config Config) (*owm.Forecast5WeatherData, error) {
	w, err := owm.NewForecast("5", "c", "en", config.APIKey)
	if err != nil {
		return nil, err
	}
	// w.DailyByName("Albuquerque", 40)
	// better use City ID to get unique responce from open weather maps
	// run below command from your linux terminal to find id of your city
	// wget -qO - http://bulk.openweathermap.org/sample/city.list.json.gz | zcat | grep -B1 -A4 Albuquerque
	// and update the config.json file with the city ID

	w.DailyByID(config.CityID, 40)
	forecast := w.ForecastWeatherJson.(*owm.Forecast5WeatherData)
	return forecast, err
}

func job(config Config) {
	w, err := getForecast5(config)
	if err != nil {
		log.Fatalln(err)
	}
	if w.City.ID != config.CityID { // open weather map did not return correct weather data, possible reason: network, server, etc. error?
		renderErrorDisp("Weather data not available.\nCheck your WiFi network!")
	} else {
		processWeatherData(&displayTxt, w)
	}

	// when we are on kidle go ahead with the display task
	if config.Kindle == 1 {
		clearDisplay()
		showImage(picFile)
	}

	// 	// for debug and test purpose show the data in a text form uncoment from => to <=
	// 	// =>
	// 	const forecastTemplate = `Weather Forecast for {{.City.Name}}:
	// {{range .List}}Date & Time: {{.DtTxt}}
	// Conditions:  {{range .Weather}}{{.Main}} {{.Description}}{{end}}
	// Temp:        {{.Main.Temp}}
	// High:        {{.Main.TempMax}}
	// Low:         {{.Main.TempMin}}

	// {{end}}
	// `

	// 	tmpl, err := template.New("forecast").Parse(forecastTemplate)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// 	if err := tmpl.Execute(os.Stdout, w); err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	// // show the downloaded forecast data in a original format
	// 	file0, err := os.Create("owm.txt")
	// 	file0.WriteString(fmt.Sprint(w))
	// 	file0.Close()

	// 	// show json nice structure
	// 	file, err := os.Create("input.txt")
	// 	if err != nil {
	// 		return
	// 	}
	// 	defer file.Close()

	// 	var jsonData string

	// 	jsonData = "["
	// 	for i, ww := range w.List {
	// 		js, _ := json.MarshalIndent(ww, "", "  ")
	// 		if i == 0 {
	// 			jsonData = jsonData + string(js)
	// 		} else {
	// 			jsonData = jsonData + ",\n" + string(js)
	// 		}
	// 	}
	// 	jsonData = jsonData + "]"
	// 	file.WriteString(jsonData)
	// 	// <==

}

func clearDisplay() {
	cmd := exec.Command("eips", "-c")
	cmd.Run()

	cmd = exec.Command("eips", "-c")
	cmd.Run()
}

func showImage(imagePath string) {
	cmd := exec.Command("eips", "-f", "-g", imagePath)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

var wg sync.WaitGroup

// global display holds data from the past today's forecasts for morning, and afternoon, when they are not a forecast any loner
var displayTxt [numDays]displayTxtType

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
	zeroDisplayTxt(&displayTxt)

	c := cron.New()
	c.AddFunc("@hourly", func() {
		job(config)
	})
	c.Start()

	job(config)

	wg.Add(1)
	wg.Wait()

	os.Exit(0)
}

func customLogf(str string, args ...interface{}) {
	fmt.Printf(str, args...)
}
