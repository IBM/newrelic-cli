#!/usr/bin/env bash

curDateFolder=$(date +%Y%m%d-%H%M%S)

basepath="./"

targetFolder="./backup_for_migration/$curDateFolder"
if [ ! -d $targetFolder ];then
    mkdir -p $targetFolder
fi

${basepath}/nr backup alertsconditions -n -s -d $targetFolder -r fail-backup-conditions.log
exitCode=$?""

if [ $exitCode == "0" ];then
    echo ""
    echo ">>>>>>>>Success, backup alertsconditions complete."
else
    echo ""
    echo ">>>>>>>>Some conditions to backup failed, begin to retry..."
fi

counter=0
while [ $exitCode != "0" ]
do
    counter=`expr $counter + 1`

    ${basepath}/nr backup alertsconditions -n -s -d $targetFolder -r fail-backup-conditions.log
    exitCode=$?""

    if [ $exitCode != "0" ];then
        echo ""
        echo "Some conditions to backup failed in this retry: "$counter"."
        if [ $counter -ge 3 ];then
            echo ""
            echo ">>>>>>>>Retry 3 times, backup quit."
            exit 1
        fi
    else
        echo ""
        echo ">>>>>>>>Success, backup alertsconditions complete."
        break
    fi    
done

echo ""

${basepath}/nr backup dashboards -s -d $targetFolder -r fail-backup-dashboards.log
exitCode=$?""

if [ $exitCode == "0" ];then
    echo ""
    echo ">>>>>>>>Success, backup dashboards complete."
else
    echo ""
    echo ">>>>>>>>Some dashboards to backup failed, begin to retry..."
fi

counter=0
while [ $exitCode != "0" ]
do
    counter=`expr $counter + 1`

    ${basepath}/nr backup dashboards -s -d $targetFolder -r fail-backup-dashboards.log
    exitCode=$?""

    if [ $exitCode != "0" ];then
        echo ""
        echo "Some dashboards to backup failed in this retry: "$counter"."
        if [ $counter -ge 3 ];then
            echo ""
            echo ">>>>>>>>Retry 3 times, backup quit."
            exit 1
        fi
    else
        echo ""
        echo ">>>>>>>>Success, backup dashboards complete."
        break
    fi    
done

echo ""
${basepath}/nr backup monitors -s -d $targetFolder -r fail-backup-monitors.log
exitCode=$?""

if [ $exitCode == "0" ];then
    echo ""
    echo ">>>>>>>>Success, backup monitors complete."
else
    echo ""
    echo ">>>>>>>>Some monitors to backup failed, begin to retry..."
fi

counter=0
while [ $exitCode != "0" ]
do
    counter=`expr $counter + 1`

    ${basepath}/nr backup monitors -d $targetFolder -r fail-backup-monitors.log
    exitCode=$?""

    if [ $exitCode != "0" ];then
        echo ""
        echo "Some monitors to backup failed in this retry: "$counter"."
        if [ $counter -ge 3 ];then
            echo ""
            echo ">>>>>>>>Retry 3 times, backup quit."
            exit 1
        fi
    else
        echo ""
        echo ">>>>>>>>Success, backup monitors complete."
	break
    fi    
done

echo ""
${basepath}/nr get users > "$targetFolder/all-in-one-bundle.user.bak"
exitCode=$?""

if [ $exitCode == "0" ];then
    echo ""
    echo ">>>>>>>>Success, backup users complete."
else
    echo ""
    echo ">>>>>>>>Backup users failed, begin to retry..."
fi

counter=0
while [ $exitCode != "0" ]
do
    counter=`expr $counter + 1`
    ${basepath}/nr get users > "$targetFolder/all-in-one-bundle.user.bak"
    exitCode=$?""

    if [ $exitCode != "0" ];then
        echo ""
        echo "Backup users failed in this retry: "$counter"."
        if [ $counter -ge 3 ];then
            echo ""
            echo ">>>>>>>>Retry 3 times, backup quit."
            exit 1
        fi
    else
        echo ""
        echo ">>>>>>>>Success, backup users complete."
        break
    fi
done
echo ""
echo "Backup for migration complete"



