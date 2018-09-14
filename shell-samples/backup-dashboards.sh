#!/usr/bin/env bash

## Can set the NEW_RELIC_APIKEY env here by comment out below export cmd
#export NEW_RELIC_APIKEY="xxx-xxxx-xxx"

curDateFolder=$(date +%Y%m%d-%H%M%S)

basepath=$(cd `dirname $0`; pwd)

targetFolder=${basepath}/backup_dashboards/$curDateFolder
if [ ! -d $targetFolder ];then
    mkdir -p $targetFolder
fi

${basepath}/nr backup dashboards -d $targetFolder -r fail-backup-dashboards.log
exitCode=$?""

if [ $exitCode == "0" ];then
    echo ""
    echo "Success, backup end."
    exit 0
else
    echo ""
    echo "Some dashboards to backup failed, begin to retry..."
fi

counter=0
while [ $exitCode != "0" ]
do
    counter=`expr $counter + 1`

    ${basepath}/nr backup dashboards -d $targetFolder -r fail-backup-dashboards.log
    exitCode=$?""

    if [ $exitCode != "0" ];then
        echo ""
        echo "Some dashboards to backup failed in this retry: "$counter"."
        if [ $counter -ge 3 ];then
            echo ""
            echo "Retry 3 times, backup end."
            exit 1
        fi
    else
        echo ""
        echo "After retry, no failed, backup end."
        exit 0
    fi    
done
