package cluster

import (
	"github.com/pattonjp/localcluster/pkg/utils"
)

func List(streamOut bool) error {
	return utils.Exec(streamOut, utils.BinK3D, "cluster", "list")
}

func SetContext(streamOut bool, name string) error {
	addArgs := []string{
		"kubeconfig", "merge", name,
		"--kubeconfig-merge-default", "--kubeconfig-switch-context",
	}

	return utils.Exec(streamOut, utils.BinK3D, addArgs...)
}

func Start(streamOut bool, name string) error {
	args := []string{"cluster", "start", name}

	return utils.Exec(streamOut, utils.BinK3D, args...)
}

func Stop(streamOut bool, name string) error {
	args := []string{"cluster", "stop", name}

	return utils.Exec(streamOut, utils.BinK3D, args...)
}

func Delete(streamOut bool, name string) error {
	addArgs := []string{"cluster", "delete", name}

	return utils.Exec(streamOut, utils.BinK3D, addArgs...)

}
