package users

import (
	"testing"

	i "github.com/IBM/newrelic-cli/test/integration"
)

func TestGetUsers(t *testing.T) {
	i.EXENRCLI("get", "users", "-o", "yaml")
}
