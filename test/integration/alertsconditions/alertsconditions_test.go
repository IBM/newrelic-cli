package alertsconditions

import (
	"testing"

	i "github.com/IBM/newrelic-cli/test/integration"
)

// func TestCreateMonitor(t *testing.T) {
// 	i.EXENRCLI("create", "monitor", "-f", "../../fixture/create/test-script-monitor.json")
// }

func TestBackupAlertsConditions(t *testing.T) {

	i.EXENRCLI("backup", "alertsconditions", "-d", "../../fixture/backup/alertsconditions", "-r", "list.log")

	i.EXEOperationCmd("rm", "-rf", "./list.log")
	i.EXEOperationCmd("rm", "-rf", "./../../fixture/backup/alertsconditions")
	i.EXEOperationCmd("mkdir", "-p", "./../../fixture/backup/alertsconditions")
}
