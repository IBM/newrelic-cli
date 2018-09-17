package alertsconditions

import (
	"testing"

	i "github.com/IBM/newrelic-cli/test/integration"
)

// func TestCreateMonitor(t *testing.T) {
// 	i.EXENRCLI("create", "monitor", "-f", "../../fixture/create/test-script-monitor.json")
// }

func TestBackupAlertsConditions(t *testing.T) {

	i.EXEOperationCmd("mkdir", "-p", "./../../fixture/output/backup/alertsconditions")

	i.EXENRCLI("backup", "alertsconditions", "-d", "../../fixture/output/backup/alertsconditions", "-r", "list.log")

	i.EXEOperationCmd("rm", "-rf", "./list.log")
	i.EXEOperationCmd("rm", "-rf", "./../../fixture/output/backup/alertsconditions")

}
