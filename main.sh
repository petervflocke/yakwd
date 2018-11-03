#!/bin/sh
/etc/init.d/framework stop
/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 1
./yakwd
/usr/bin/lipc-set-prop -- com.lab126.powerd preventScreenSaver 0
/etc/init.d/framework start