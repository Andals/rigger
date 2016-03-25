#!/bin/bash

cur_path=`dirname $0`

cd $cur_path/../
rigger_home=`pwd`
cd ../
prj_home=`pwd`
cd - > /dev/null

/usr/local/php/bin/php $rigger_home/server/run.php $* rigger_conf=$rigger_home/client/rigger_conf.php prj_home=$prj_home
