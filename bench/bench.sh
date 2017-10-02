#!/usr/bin/env bash

#base=${1:-1}
#echo $((`git rev-list --all|wc -l` + $base))
file='bench_'
file=$file`git describe --always`_`date +%Y%m%d_%T`
#echo $file
go test -bench=. > $file
