#!/usr/bin/env bash

usage() {
	echo "usage:"
	echo "  wtf collector [error|warning|ok] test participant counter msg"
}

report() {
	redis="redis-cli -h $collector"
	t=`$redis --raw time`
	t0=`sed -n 1p <<< "$t"`
	t1=`sed -n 2p <<< "$t"`
	#echo "[report] time $t0 $t1"
	
	
	key=`printf "%s:%s:%s" $test_ $participant $counter`
	value=`printf "%s:::%s" $level $msg`

	#echo "[report] $key"
	#echo "[report] $value"

	$redis SET $key $value &> /dev/null

	time_key=`printf "%s:~time~" $key`
	$redis DEL $time_key &> /dev/null
	$redis RPUSH $time_key $t0 &> /dev/null
	$redis RPUSH $time_key $t1 &> /dev/null
}

if [[ "$#" -ne 6 ]]; then
	usage
	exit 1
fi

collector=$1
level=$2
test_=$3
participant=$4
counter=$5
msg=$6

case $level in
	error|warning|ok)
		report
		;;
	*)
		usage
		;;
esac


