#!/usr/bin/env bash

usage() {
	echo "usage:"
	echo "  wtf [error|warning|ok] test participant counter msg"
}

report() {
	t=`redis-cli --raw time`
	t0=`sed -n 1p <<< "$t"`
	t1=`sed -n 2p <<< "$t"`
	echo "[report] time $t0 $t1"
	
	
	key=`printf "%s:%s:%s" $test_ $participant $counter`
	value=`printf "%s:::%s" $level $msg`

	echo "[report] $key"
	echo "[report] $value"

	redis-cli SET $key $value

	time_key=`printf "%s:~time~" $key`
	redis-cli DEL $time_key
	redis-cli RPUSH $time_key $t0
	redis-cli RPUSH $time_key $t1
}

if [[ "$#" -ne 5 ]]; then
	usage
	exit 1
fi

level=$1
test_=$2
participant=$3
counter=$4
msg=$5

case $level in
	error|warning|ok)
		report
		;;
	*)
		usage
		;;
esac


