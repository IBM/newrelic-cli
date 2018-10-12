# newrelic-cli

[![Build Status](https://travis-ci.com/IBM/newrelic-cli.svg?branch=master)](https://travis-ci.com/IBM/newrelic-cli)

New Relic CLI is a command line tool which is used to operate New Relic objects(Synthetic monitors, alert policies, conditions, account users etc). You can use it easily to list/get/create/delete these objects. It can be used to backup your New Relic configuration data and restore in the future. It is easy to be used other than calling different REST API endpoints.

## Command

Command | Subcommand | Resource | arguments | flag
-- | -- | -- | -- | --
nr | get | users | - | 
nr | get | user | &lt;id&gt; | 
nr | get | monitors | - | 
nr | get | monitor | &lt;id&gt; | 
nr | get | labels | - | 
nr | get | labelsmonitors | &lt;category:label&gt; | 
nr | get | alertspolicies | - | 
nr | get | alertsconditions | - | 
nr | get | alertschannels | - | 
nr | get | dashboards | - | 
nr | get | dashboard | &lt;id&gt; | 
nr | create | monitor | - | -f &lt;monitor_sample.json&gt;
nr | create | alertspolicies | - | -f &lt;alertspolicies_sample.json&gt;
nr | create | alertsconditions | - | -f &lt;alertsconditions_sample.json&gt;
nr | create | alertschannels | - | -f &lt;alertschannels_sample.json&gt;
nr | add | alertschannels | &lt;id&gt; &lt;category:label&gt; | 
nr | update | monitor | - | -f &lt;monitor_sample.json&gt;
nr | update | alertspolicies | - | -f &lt;alertspolicies_sample.json&gt;
nr | update | alertsconditions | - | -f &lt;alertsconditions_sample.json&gt;
nr | update | alertschannels | - | -f &lt;alertschannels_sample.json&gt;
nr | patch | monitor | - | -f &lt;monitor_sample.json&gt;
nr | delete | monitor | &lt;id&gt; | 
nr | delete | alertspolicies | &lt;id&gt; | 
nr | delete | alertsconditions | &lt;id&gt; | 
nr | delete | alertschannels | &lt;id&gt; | 
nr | delete | labelsmonitors | &lt;id&gt; &lt;category:label&gt; | 
nr | insert | customevents | - | -f &lt;custom_events.json&gt;<br> -i &lt;New Relic insert key&gt;<br> -a &lt;New Relic account ID&gt;<br>
nr | backup | monitors | - | -d &lt;backup_folder&gt;<br> -r &lt;result_file.log&gt;<br>
nr | backup | alertsconditions | - | -d &lt;backup_folder&gt;<br> -r &lt;result_file.log&gt;<br>
nr | backup | dashboards | - | -d &lt;backup_folder&gt;<br> -r &lt;result_file.log&gt;<br>
nr | restore | monitors | - | -d &lt;monitors_folder&gt;<br> -f &lt;monitor_filenames&gt;<br> -F &lt;file_contains_names&gt;<br> -m [skip\\|override\\|clean]<br> -r &lt;result_file.log&gt;<br>
nr | restore | alertsconditions | - |  -d &lt;alertsconditions_folder&gt;<br> -f &lt;alertscondition_filenames&gt;<br> -F &lt;file_contains_names&gt;<br> -m [skip\\|override\\|clean]<br> -r &lt;result_file.log&gt;<br>
nr | restore | dashboards | - |  -d &lt;dashboards_folder&gt;<br> -f &lt;dashboard_filenames&gt;<br> -F &lt;file_contains_names&gt;<br> -m [skip\\|override\\|clean]<br> -r &lt;result_file.log&gt;<br>

## To start using nr CLI

### Getting Started with the nr CLI (A quick sample to get all users in New Relic account)

* __Set environment variable__ `NEW_RELIC_APIKEY`

Define New Relic admin API key in environment by `export` cmd on Linux OS like this:
<br>
`export NEW_RELIC_APIKEY=xxxx-xxxxxxx-xxxxx-xxxxxx`
<br><br>
`xxxx-xxxxxxx-xxxxx-xxxxxx` is the New Relic admin API key. How to find the admin API key in your New Relic account, reference this doc, __Activate Admin user's API key__: [REST API keys](https://docs.newrelic.com/docs/apis/getting-started/intro-apis/access-rest-api-keys)<br><br>



* __Get all users info in current New Relic account__

Use __nr get users__ command like this:<br>
```
$ nr get users
ID        FirstName   LastName    Email                Role
2071178   Tom       Smith       xxx@test.com    admin
2000900   Jack        Xi        xxx@test.com     admin
```


<br>Define the output format as JSON using `-o json` argument<br>
```
$ nr get users -o json
{
  "users": [
    {
      "id": 2071178,
      "first_name": "Tom",
      "last_name": "Smith",
      "email": "xxx@tom.com",
      "role": "admin"
    },
......
```


<br>Define the output format as JSON using `-o YAML` argument<br>
```
$ nr get users -o yaml
users:
- email: xxx@test.com
  first_name: Tom
  id: 2071178
  last_name: Smith
  role: admin
......
```

* __Use proxy__

Can configure proxy if the target machine can not directly connect to newrelic.com

`export NEW_RELIC_PROXY=http://<user>:<password>@<ip>:<port>`

Like:<br>
`export NEW_RELIC_PROXY=http://user1:password1@9.42.95.127:3128`


* __Configure retries__

Can configure retries for calling NewRelic REST API while some network issues, all NewRelic REST API callings in CLI would follow the retries. The default retries value is __3__ if `RETRIES` not configure 

`export RETRIES=<times>`

Like:<br>
`export RETRIES=5`

* __Return codes__

The nr CLI uses exit codes, which help with scripting and confirming that a command has run successfully. For example, after you run a nr CLI command, you can retrieve its return code by running echo $? (on Windows, echo %ERRORLEVEL%). If the return code is 0, the command was successful.

__A sample to use `return code` in shell scripts for CI/CD pipeline__
To restore monitors, we check if successful by the `return code`, if failed, it would output the monitor file names to the log file `fail-restore-monitors.log`(the file name you can customize by `-r` argument), we continue to retry to restore all monitors that failed, the monitor names were stored in `fail-restore-monitors.log` and we use `-F` argument to tell `nr` what monitors we want to retry, we also add a counter, if retry times exceed 3, it would end.

__shell scripts to restore monitors with retry:__
```
#!/bin/bash

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
```


* __nr help command__

The nr help command lists the nr CLI commands and a brief description of each. Passing the -h flag to any command lists detailed help, including any aliases. For example, to see detailed help for `nr get`, run:

```
$ nr get -h
Display one or many NewRelic resources.

Usage:
  nr get [command]

Available Commands:
  alertschannels   Display all alerts_channels.
  alertsconditions Display alert conditions by alert id.
  alertspolicies   Display all alerts_policies.
  labels           Display all labels.
  labelsmonitors   Display monitros by label.
  monitor          Display a single monitor by id.
  monitors         Display all synthetics monitors.
  user             Display a single user by id.
  users            Display all users.

Flags:
  -h, --help                    help for get
  -o, --output string           Output format. table/json/yaml are supported (default "table")
  -t, --type-condition string   Alert condition type. Only used for 'alertsconditions' command. all|conditions|synthetics|ext|plugin|nrql are supported (default "all")
```

<br>for `nr get users`, run:

```
$ nr get users -h
Display all users.

Usage:
  nr get users [flags]

Examples:
* nr get users
* nr get users -o json
* nr get users -o yaml
* nr get users -i 2102902
* nr get users -i 2102902,+801314

Flags:
  -e, --email string   email to filter returned result. can't specify emails
  -h, --help           help for users
  -i, --id string      user id(s) to filter returned result. use ',+' to separate ids

Global Flags:
  -o, --output string           Output format. table/json/yaml are supported (default "table")
  -t, --type-condition string   Alert condition type. Only used for 'alertsconditions' command. all|conditions|synthetics|ext|plugin|nrql are supported (default "all")
```


## To start developing nr CLI

### Prerequisite

* Golang 1.9 or 1.9+

* Golang `dep`
<br>Install `dep`:<br>
`go get -u github.com/golang/dep/cmd/dep`


### Build newrelic-cli project

* git clone `newrelic-cli`
* Enter project root folder
* Run `make deps`
* Run `make build`


### Cross compilation by using gox

* Install `gox`
* Enter project root folder
* Run `gox -os "windows linux darwin" -arch "amd64"`


### How to run unit test
* `export NEW_RELIC_APIKEY=<Your NewRelic API Key>` 
* `make test`

## Changelog
[Changelog](https://github.com/IBM/newrelic-cli/blob/master/CHANGELOG.md)

## Acknowledgement

Special thanks to [Huang Wei](https://github.com/Huang-Wei) who proposed this good idea and developed the initial version.