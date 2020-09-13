package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

func main() {
	if test() == nil {
		build("windows", "amd64")
		build("linux", "amd64")
		build("linux", "arm")
		build("linux", "arm64")
	}
}

func test() error {
	return execute(
		[]string{"go", "test", "./..."})
}

func build(targetOs string, targetArch string) error {
	targetDir := path.Join("bin", fmt.Sprintf("%s-%s", targetOs, targetArch))
	os.MkdirAll(targetDir, 0777)

	targetDef := []string{
		fmt.Sprintf("GOOS=%s", targetOs),
		fmt.Sprintf("GOARCH=%s", targetArch)}
	if targetArch == "arm" {
		targetDef = append(targetDef, "GOARM=7")
	}

	return execute(
		[]string{
			"go",
			"build",
			"-a",
			"-o",
			targetDir,
			"github.com/soerenkoehler/chdiff-go/chdiff"},
		targetDef...)
}

func execute(cmdline []string, env ...string) error {
	fmt.Println(cmdline, env)

	proc := exec.Command(cmdline[0], cmdline[1:]...)
	proc.Env = append(os.Environ(), env...)
	pipeOut, _ := proc.StdoutPipe()
	pipeErr, _ := proc.StderrPipe()

	output := make(chan string)
	defer close(output)

	go copyOutput(pipeOut, output)
	go copyOutput(pipeErr, output)
	go printOutput(output)

	err := proc.Run()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	return err
}

func copyOutput(src io.Reader, dest chan<- string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		dest <- scanner.Text()
	}
}

func printOutput(src <-chan string) {
	for line := range src {
		fmt.Println(line)
	}
}
