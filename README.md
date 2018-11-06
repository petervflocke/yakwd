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

+ No real jailbreak is needed, just a root access via ssh.
  + Here how to do this => https://github.com/petervflocke/yakwd/blob/master/enablesshoverwifi.md
+ [Screenshot](https://github.com/petervflocke/yakwd#screenshots)

### Remarks:

* Inspired by https://github.com/DDRBoxman/kindle-weather
* Tested on Kindle 4 Non Touch https://wiki.mobileread.com/wiki/Kindle4NTHacking
* To run it you need an API key from https://openweathermap.org/api
* If you do not want to install Golang, use a ready binary from bin folder. It was tested on **Kindle 4 NT**


### Quick Install:

1. Download or clone this github repository 
2. Copy from project bin folder file `yakwd` to your kindle (`/var/tmp/root/`)
3. Copy the `font` folder to to your kindle (`/var/tmp/root/`)
4. Copy `yakwd.json.sample` to to your kindle (`/var/tmp/root/yakwd.json`) - update it as follow
  * your own open weather API key. 
  * the city ids to be added to the table (at least one) 
  * an index for the first city to be displayed. Index starts from 0
5. Copy `yakwd.sh` to your kindle `/var/tmp/root/`
6. Run from ssh terminal on kindle: `./yakwd.sh`

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
* Rename `yakwd.json.sample` to `yakwd.json` - update as above in the quick installation section
* Change to `"kindle"      : 0,` to bypass all kindle dependent functions
* Run go source from yakwd folder: `go run yakwd`

The ouput graphic file `out.png` is saved `/tmp/`  folder.


****

### How to pernamently install on Kindle:

*On Kindle:*
1. Login via ssh and change root system to writeable: `mntroot rw`

*On PC:* from the project folder: (kindle = ip address of the kindle)
1. Prepare yakwd.jason (API-Keys, City IDs, City index) see above descriptions
2. scp bin/yakwd root@kindle:/usr/local/bin/
3. scp fonts/kindleweathersr.ttf root@kindle:/usr/java/lib/fonts/
4. scp yakwd.json root@kindle:/etc/
5. scp yakwd.sh root@kindle:/usr/local/sbin/

*On Kindle:*
1. Login via ssh and change root system to writeable: `mntroot ro`

Done!

To start the application (no launcher on kindle) run from your PC or mobile:<BR>
`ssh root@kindle "/usr/local/sbin/yakwd.sh >/dev/null 2>&1 &"`

****

### ToDo

* Add init script for start at boot
* Reduce current consupmtion by adding sleep mode or switching off wifi when not needed
* Add family calendar view option, warnings from home automations system ,etc
* Out of scope: creating a Kindle "update" package, manual installation will be always required


### "ScreenShots"

Example from the alfa v1.0 Version with an own footer of location, time and battery<br>
Pictures have a spot in the upper corner, it's a broken screen :(

![live foto](https://github.com/petervflocke/yakwd/blob/master/Docs/kindle-live.jpg)