# yakwd #

(Alfa version of a) stand alone app for weather forecast direct download from https://openweathermap.org and display on a Kindle 4 NT<br>

### Main Features
+ No additional server nor local kindle python is needed.
+ No additional graphic software converter (e.g. imagemagick, pngcrush, etc) is needed. Go code produces 8 bit gray scale pics, eips ready entirely on the Kindle.
+ App displays forecast for 3 days today, tomorrow and a day after tomorrow. For each day 3 time zones are reflected (morning, afternoon and evening). Weather picture is updated at every full hour
+ Kindle 4 keyboard functions:
  + _Back_ : exits the application
  + _Virtsual Keyboard_ : refresh the forecast now
  + _Right Next Page_ : switch to the next city on the list
  + _Left Next Page_ : switch to the previous city on the list

### Remarks:

* Inspired by https://github.com/DDRBoxman/kindle-weather
* Tested on Kindle 4 No Touch https://wiki.mobileread.com/wiki/Kindle4NTHacking
* To run it you need an API key from https://openweathermap.org/api
* If you do not want to install Golang, use a ready binary from bin folder. It was tested on Kindle 4 NT
* Full starting script for the init process will follow soon

### Quick Install:

1) copy from bin project folder file "yakwd" to your kindle (/var/tmp/root/)
2) copy to the folder from bin folder file `yakwd`
* the `font` folder
* copy `yakwd.json.sample` to `yakwd.json` - update it as follow
  * your own open weather API key. 
  * the city ids to be added to the table (at least one) 
  * an index for the first city to be displayed. Index starts from 0
* main.sh 
run from ssh terminal: `./yakwd.sh`

### Long Install and some config options

* clone git repository
* install all used packages
* compile go source, see build file for target kindle architecture
* copy results including config, scripts and fonts files to your kindle.

In the yakwd.json you can configure which fonts shall be used. These are the ones I use:
```json
  "iconfont"    : "/usr/java/lib/fonts/kindleweathersr.ttf",
  "txtfont"     : "/usr/java/lib/fonts/Helvetica_LT_65_Medium.ttf" 
```
For the test and quick installation in a /var/tmp/root/ the local folder `./font/` with Robotosr.ttf and kindleweathersr.ttf can be used. If the fonts are not found, yakwd tries to load them from local .font folder.

For the default instalation kindleweathersr.ttf has to be copied into the kindle font directory. <BR>
1. on kindle Make the root rw (`mntroot rw`)
2. from pc (s)copy `./font/kindleweathersr.ttf` to `kindle_ip://usr/java/lib/fonts/`

### Testing on PC without Kindle
* Clone git repository
* Install all used packages
  * `go get github.com/robfig/cron`
  * `go get github.com/briandowns/openweathermap`
  * `go get github.com/fogleman/gg`
  * `go get github.com/fogleman/gg`
* Rename `yakwd.json.sample` to `yakwd.json` - update as above in the quick instllation section
* Change to `"kindle"      : 0,` to bypass all kindle dependent functions
* Run go source from yakwd folder: `go run yakwd`

The ouput graphic file `out.png` is saved `/tmp/`  folder.


****

Description under development (do not use for now, I'm documenting my own steps):
*On Kindle:*
+ ~~mntroot rw~~

*On PC:"
1. ~~scp yakwd root@192.168.0.144:/usr/local/bin/~~
2. ~~scp fonts/kindleweathersr.ttf root@192.168.0.144:/usr/java/lib/fonts/~~
3. ~~scp yakwd.json root@192.168.0.144:/etc/~~
4. ~~scp yakwd_init.sh root@192.168.0.144:/etc/init.d/~~

****


### ToDo

* Add init script for start at boot
* Reduce current consupmtion by adding sleep mode or switching off wifi when not needed
* Add calendar view option
* Out of scope: Creating a Kindle "update" package, manul instllation will bbe required


### "ScreenShots"

Example from the alfa v1.0 Version with an own footer of location, time and battery<br>
Pictures have a spot in the upper corner, it's a broken screen :(

![live foto](https://github.com/petervflocke/yakwd/blob/master/Docs/kindle-live.jpg)
<br>

Earlier version, still with Kindle bar at the top.
![live foto](https://github.com/petervflocke/yakwd/blob/master/Docs/kindle-live-2.jpg)
