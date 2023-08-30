package utils

import (
	"fmt"
	"strings"
	"time"
)

func WaitForDeployment(app, ns, label string) error {
	o, err := execWait(BinKubectl,
		"--namespace", ns,
		"get", "pod",
		"-l", fmt.Sprintf("%s=%s", label, app),
		"--output", "jsonpath='{.items[*].status.phase}'",
	)
	if err != nil {
		return err
	}

	if !strings.Contains(o, "Running") {
		fmt.Printf("waiting for %s\n", app)
		time.Sleep(1 * time.Second)
		WaitForDeployment(app, ns, label)
	}
	return nil
}
