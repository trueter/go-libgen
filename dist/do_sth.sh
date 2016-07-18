#!/bin/bash


components='webout.sh webin.go'

function setup {
	echo "Running setup"
	echo "creatng dir structure"
	wd=`pwd`
	mkdir fdb && cd fdb && mkdir in out && cd in && mkdir tomobi toepub topdf totxt 
	cd $wd
}


function usage {
	echo some helpfull message
	echo '[up|down|setup|help]'
	exit
}

function up {
	nohub ./foo &
	nohub ./bar &
	exit
}

function down {
	killall foooo
	killall baaaar
}


function chkdp {
# does ebook_convert and redis exist
# is redis running? Can we connect?
# and nohub? and the dirs? and killall and what not
}

if [ -z ${1+x} ]; then
	usage
elif [ "$1" == 'setup' ] ; then
	setup
elif [ "$1" == 'up' ] ; then
	up
elif [ "$1" == 'down' ] ; then
	stop 
else    # elif  [[ "$1" == "--help" || "$1" == "-h" ]]; then
	usage
fi



