#!/bin/sh

_FUNCTIONS=/etc/rc.d/functions
[ -f ${_FUNCTIONS} ] && . ${_FUNCTIONS}

YAKWD=/usr/bin/yakwd
YAKWDBIN=yakwd

[ ! -x $YAKWD ] && exit 1

run() {
	/etc/init.d/framework stop
	/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 1
	/usr/sbin/eips -s 5
	$YAKWD 
	/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 0
	/etc/init.d/framework start	
}


start_YAKWD()
{
	msg "starting yakwd app" I
	run &
}

stop_YAKWD()
{
	killall -9 $YAKWDBIN
}

case "$1" in

    start)
    	start_YAKWD
	;;

    stop)
    	stop_YAKWD
	;;
	
    restart)
    	echo "not implemented"
    ;;
	
    status)
    	echo "not implemented"
    ;;

    *)
	msg "Usage: $0 start" W
	exit 1
	;;
esac

exit 0
