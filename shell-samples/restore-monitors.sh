#!/usr/bin/env bash

noNRKey="No NEW_RELIC_APIKEY detected."
noNRKey=${noNRKey}"\n\n"
noNRKey=${noNRKey}"Please export New Relic API key."
noNRKey=${noNRKey}"\n\n"
noNRKey=${noNRKey}"Example:\n"
noNRKey=${noNRKey}"  export NEW_RELIC_APIKEY=xxx-xxxxx-xx"
noNRKey=${noNRKey}"\n"
if [ $NEW_RELIC_APIKEY"" == "" ];then
    echo -e "${noNRKey}"
    exit 1
fi

basepath=$(cd `dirname $0`; pwd)

${basepath}/nr restore monitors -d ${basepath}/backup_monitors_folder -r fail-restore-monitors.log
exitCode=$?""

if [ $exitCode == "0" ];then
    echo ""
    echo "Success, restore end."
    exit 0
else
    echo ""
    echo "Some monitors to restore failed, begin to retry..."
fi

counter=0
while [ $exitCode != "0" ]
do
    counter=`expr $counter + 1`

    ${basepath}/nr restore monitors -F ${basepath}/fail-restore-monitors.log -r fail-restore-monitors.log
    exitCode=$?""

    if [ $exitCode != "0" ];then
        echo ""
        echo "Some monitors to restore failed in this retry: "$counter"."
        if [ $counter -ge 3 ];then
            echo ""
            echo "Retry 3 times, restore end."
            exit 1
        fi
    else
        echo ""
        echo "After retry, no failed, restore end."
        exit 0
    fi    
done
