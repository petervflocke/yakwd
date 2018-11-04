# yakwd

(Alfa version of a) stand alone app for weather forecast direct download from https://openweathermap.org and display on a Kindle 4 NT<br>
<br>No additional server needed.<br>
<br>No additional graphic converter (e.g. imagemagick, pngcrush, etc) is needed. Go code produces 8 bit gray scale pics, eips ready.<br>
<br>App displays forecast for 3 days today, tomorrow and a day after tomorrow. For each day 3 time zones are reflected (morning, afternoon and evening)<br>
<br>Early but runnable code :)<br>

Remarks:
---
* If you do not want to install Golang, use a ready binary from bin folder. It was tested on Kindle 4 NT
* Starting script for the init process will follow soon
* Remark: Pictures have a spot in the upper corner, it is a broken screen.

Install:
---
1) copy from bin file "yakwd" to your kindle#
2) copy to the folder with yakwd 
* font folder
* config.json.sample and rename to config.json - update it with the city id and your open weather API key
* main.sh 

run main.sh

---

Example from the alfa v1.0 Version with an own footer of location, time and battery

![live foto](https://github.com/petervflocke/yakwd/blob/master/Docs/kindle-live.jpg)
<br>

Earlier version, still with Kindle bar at the top.
![live foto](https://github.com/petervflocke/yakwd/blob/master/Docs/kindle-live-2.jpg)
