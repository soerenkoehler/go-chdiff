package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var pathCWD, _ = os.Getwd()
var envGOBIN = fmt.Sprintf("GOBIN=%s", path.Join(pathCWD, "bin"))
var cmdINSTALL = []string{"go", "install", "github.com/soerenkoehler/chdiff-go/chdiff"}

func main() {
	execute(cmdINSTALL, envGOBIN)
}

func execute(cmdline []string, env ...string) {
	proc := exec.Command(cmdline[0], cmdline[1:]...)

	proc.Env = append(os.Environ(), env...)
	stdout, _ := proc.StdoutPipe()
	stderr, _ := proc.StderrPipe()

	proc.Start()
	stdoutContent, _ := ioutil.ReadAll(stdout)
	stderrContent, _ := ioutil.ReadAll(stderr)

	if err := proc.Wait(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%s%s", stdoutContent, stderrContent)
}
