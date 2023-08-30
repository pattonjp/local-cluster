package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type Dep struct {
	Name    string
	Version string
}

var (
	BinK3D     = Dep{"k3d", "latest"}
	BinKubectl = Dep{"kubectl", "latest"}
	BinHelm    = Dep{"helm", "latest"}
	BinMkCert  = Dep{"mkcert", "latest"}
	asdf       = Dep{Name: "asdf"}
)

var allDeps = []Dep{
	BinMkCert,
	BinK3D,
	BinKubectl,
	BinHelm,
}

func InstallAllDeps() error {
	for _, d := range allDeps {
		install(d)
	}
	return nil
}

func MustCheckAllDeps() {
	for _, d := range allDeps {
		MustCheckDep(d)
	}
}

func CheckDep(app Dep) bool {
	_, err := exec.LookPath(app.Name)
	return err == nil
}

func MustCheckDep(app Dep) {
	if !CheckDep(app) {
		panic(fmt.Sprintf("dependency was not found %s. please install\n", app))
	}
}

func scanner(reader io.ReadCloser) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println(scanner.Text())

	}
}
func Exec(stream bool, name Dep, args ...string) error {
	if stream {
		return execSteam(name, args...)
	} else {
		_, err := execWait(name, args...)
		return err
	}
}
func execWait(dep Dep, args ...string) (string, error) {
	cmd := exec.Command(dep.Name, args...)
	out, err := cmd.Output()
	return string(out), err

}

func execSteam(dep Dep, args ...string) error {
	cmd := exec.Command(dep.Name, args...)
	var stdout io.ReadCloser
	var stdErr io.ReadCloser
	var err error
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return err
	}
	if stdErr, err = cmd.StderrPipe(); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	go scanner(stdout)
	go scanner(stdErr)

	if err := cmd.Wait(); err != nil {
		return err
	}
	return err
}
