package monitors

import (
	"testing"

	i "github.com/IBM/newrelic-cli/test/integration"
)

// func TestCreateMonitor(t *testing.T) {
// 	i.EXENRCLI("create", "monitor", "-f", "../../fixture/create/test-script-monitor.json")
// }

func TestBackupMonitors(t *testing.T) {

	i.EXEOperationCmd("mkdir", "-p", "./../../fixture/output/backup/monitors")

	i.EXENRCLI("backup", "monitors", "-d", "../../fixture/output/backup/monitors", "-r", "list.log")

	i.EXEOperationCmd("rm", "-rf", "./list.log")
	i.EXEOperationCmd("rm", "-rf", "./../../fixture/output/backup/monitors")

}
