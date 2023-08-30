package cluster

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/pattonjp/localcluster/pkg/utils"
)

type args struct {
	key string
	val string
}

type stage int

const (
	initial stage = iota
	base
	app
)

var (
	stageNames = [...]string{
		"initial",
		"base",
		"app",
	}
)

type HelmDeployOptions struct {
	Name      string `yaml:"name"`
	Domain    string `yaml:"-"`
	Stage     stage  `yaml:"stage"`
	Namespace string `yaml:"namespace"`
	Chart     string `yaml:"chart"`
	RepoName  string `yaml:"repoName"`
	RepoURL   string `yaml:"repoURL"`
	helmArgs  []args `yaml:"-"`
	AddCerts  bool   `yaml:"createCerts"`
	Version   string `yaml:"version"`
}

func (stg stage) String() string {

	if int(stg) >= len(stageNames) {
		return ""
	}
	return stageNames[stg]

}

func (stg *stage) MarshalText() ([]byte, error) {
	return []byte(stg.String()), nil
}
func (stg *stage) UnmarshalText(b []byte) error {
	str := string(b)
	for i, name := range stageNames {
		if name == str {
			*stg = stage(i)
			return nil
		}
	}
	*stg = app
	return nil
}

func Setup(stream bool, config Config) error {
	list := map[string]string{}
	for _, ch := range config.Charts {
		list[ch.RepoName] = ch.RepoURL

	}
	for name, url := range list {
		utils.Exec(stream, utils.BinHelm, "repo", "add", name, url)
		utils.Exec(stream, utils.BinHelm, "repo", "update", name)
	}
	return nil
}

func DeployChart(name string, config Config) error {
	chart := config.GetChart(name)
	if chart == nil {
		return errors.New("could not find chart")
	}
	return chart.deploy(config)
}

func AvailableCharts(config Config) error {
	pKV := func(key string, val string) {
		fmt.Printf("%s: %s\n", key, val)
	}
	fmt.Println("-------------------------")
	fmt.Println("Available charts")
	for _, opts := range config.Charts {
		fmt.Println("-------------------------")
		pKV("Name", opts.Name)
		pKV("Stage", opts.Stage.String())
		pKV("Namespace", opts.Namespace)
		pKV("Chart", opts.Chart)
		pKV("Repo", opts.RepoName)
		pKV("Url", opts.RepoURL)

	}
	fmt.Println("-------------------------")
	return nil
}

func Deploy(config Config) error {
	utils.WaitForDeployment("kube-dns", "kube-system", "k8s-app")
	utils.WaitForDeployment("metrics-server", "kube-system", "k8s-app")

	for _, hdo := range config.Charts {
		fmt.Println("deploying ", hdo.Chart)
		if err := hdo.deploy(config); err != nil {
			return err
		}
	}

	return nil
}

func (do *HelmDeployOptions) deploy(config Config) error {
	if do == nil {
		panic("helm options nil")
	}
	utils.EnsureNamespace(do.Namespace)
	do.Domain = config.Domain
	valFilePaths, err := GetValuesFilesFor(do.Name)
	if err != nil {
		return err
	}
	valuesFiles := []string{}
	for _, vf := range valFilePaths {
		f, err := getAssetPath(vf, do)
		if err != nil {
			return err
		}
		valuesFiles = append(valuesFiles, f)
	}

	if do.AddCerts {
		utils.CreateLocalCert(config.Domain)
		if err := utils.AddCerts(do.Namespace); err != nil {
			return err
		}
	}

	args := []string{
		"upgrade",
		"--install",
		"--create-namespace",
		"--wait",
		"--timeout", "300m",
		"-n", do.Namespace,
	}

	if do.Version != "" {
		args = append(args, "--version", do.Version)
	}
	for _, arg := range do.helmArgs {
		args = append(args, fmt.Sprintf("--set %s=%s", arg.key, arg.val))
	}

	args = append(args, do.Name, do.Chart)
	for _, v := range valuesFiles {
		args = append(args, "-f", v)
	}
	if err := utils.Exec(true, utils.BinHelm, args...); err != nil {
		return err
	}

	kcFiles, err := GetNonValuesFilesFor(do.Name)
	if err != nil {
		return err
	}
	for _, f := range kcFiles {
		lf, err := getAssetPath(f, do)
		if err != nil {
			return err
		}
		err = utils.Exec(true, utils.BinKubectl, "apply", "-n", do.Namespace, "-f", lf)
		if err != nil {
			return err
		}
	}
	return nil

}
