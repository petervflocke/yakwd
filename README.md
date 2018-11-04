# yakwd

(Alfa version of a) stand alone app for weather forecast direct download from https://openweathermap.org and display on a Kindle 4 NT<br>
<br>No additional server needed.<br>
<br>No additional graphic converter (e.g. imagemagick, pngcrush, etc) is needed. Go code produces 8 bit gray scale pics, eips ready.<br>
<br>App displays forecast for 3 days today, tomorrow and a day after tomorrow. For each day 3 time zones are reflected (morning, afternoon and evening)<br>
<br>Early but runnable code :)<br>

Remarks:
---
* Inspired by https://github.com/DDRBoxman/kindle-weather
* Tested on Kindle 4 No Touch https://wiki.mobileread.com/wiki/Kindle4NTHacking
* To run it you need an API key from https://openweathermap.org/api
* If you do not want to install Golang, use a ready binary from bin folder. It was tested on Kindle 4 NT
* Full starting script for the init process will follow soon

Quick Install:
---
1) copy from bin project folder file "yakwd" to your kindle
2) copy to the folder with yakwd 
* font folder
* config.json.sample and rename to config.json - update it with<br>
+ your own open weather API key. 
+ the city ids to be added to the table (at least one) 
+ an index for the firts city to be displayed. Index starts from 0
* main.sh 

run main.sh

Long Install
---
* clone git repository
* install all used packages
* compile go source, see build file for target kindle architecture
* copy results including config, scripts and fonts files to your kindle.


Example from the alfa v1.0 Version with an own footer of location, time and battery<br>
Pictures have a spot in the upper corner, it's a broken screen :(

![live foto](https://github.com/petervflocke/yakwd/blob/master/Docs/kindle-live.jpg)
<br>

Earlier version, still with Kindle bar at the top.
![live foto](https://github.com/petervflocke/yakwd/blob/master/Docs/kindle-live-2.jpg)
