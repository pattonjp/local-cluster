package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func CertInit() {
	install(BinMkCert)
	execSteam(BinMkCert, "install")
}

func CreateLocalCert(domain string) {
	if _, err := os.Stat("certs"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("certs", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	args := []string{
		"-cert-file",
		"certs/local-cert.pem",
		"-key-file",
		"certs/local-key.pem",
		domain,
		fmt.Sprintf("*.%s", domain),
		"localhost",
		"127.0.0.1",
		"::1",
	}
	execSteam(BinMkCert, args...)
}

func EnsureNamespace(ns string) error {
	str, err := execWait(BinKubectl,
		"create", "namespace", ns,
	)
	if err != nil {
		fmt.Println(str)
	}
	return nil

}

func AddCerts(ns string) error {
	_, err := execWait(BinKubectl,
		"-n", ns, "delete", "secret", "tls-secret", "--ignore-not-found=true",
	)
	if err != nil {
		return err
	}

	_, err = execWait(BinKubectl,
		"-n", ns, "create", "secret", "tls", "tls-secret",
		"--cert", "certs/local-cert.pem",
		"--key", "certs/local-key.pem",
	)
	if err != nil {
		return err
	}

	return nil
}
