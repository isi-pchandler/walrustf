#!/usr/bin/env bash

usage() {
	echo "usage:"
	echo "  wtf collector [error|warning|ok] test participant msg"
}

report() {
	redis="redis-cli -h $collector"
	t=`$redis --raw time`
	t0=`sed -n 1p <<< "$t"`
	t1=`sed -n 2p <<< "$t"`
	
	key=`printf "%s:%s:%s:%s" $test_ $participant $t0 $t1`
	value=`printf "%s:::%s" $level $msg`

	$redis SET $key $value &> /dev/null
}

if [[ "$#" -ne 5 ]]; then
	usage
	exit 1
fi

collector=$1
level=$2
test_=$3
participant=$4
msg=$5

case $level in
	error|warning|ok)
		report
		;;
	*)
		usage
		;;
esac


