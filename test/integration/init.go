package integration

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

var (
	NEW_RELIC_APIKEY string
)

func init() {
	NEW_RELIC_APIKEY = ""
}

func EXENRCLI(args ...string) {

	cmd := exec.Command("../../../nr", args...)

	// env := os.Environ()
	// env = append(env, "NEW_RELIC_APIKEY="+NEW_RELIC_APIKEY)
	// cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Printf("output for failure: %q\n", out.String())
		fmt.Println()
		log.Fatal(err)
	}
	fmt.Printf("output: %q\n", out.String())

}

func EXEOperationCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("output for failure: %q\n", out.String())
		fmt.Println()
		log.Fatal(err)
	}
	fmt.Printf("output: %q\n", out.String())
}
