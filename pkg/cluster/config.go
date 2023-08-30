package cluster

import (
	"mutterio/localdev/pkg/utils"
	"os"
	"path/filepath"
	"sort"

	_ "embed"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

const confName = ".localcluster"

//go:embed charts.yaml
var chartByte []byte

type Config struct {
	Name           string              `yaml:"name"`
	Domain         string              `yaml:"domain"`
	AgentCount     int                 `yaml:"agents"`
	ServerCount    int                 `yaml:"servers"`
	ServerPorts    []int               `yaml:"ports"`
	ChartNames     []string            `yaml:"include"`
	Charts         []HelmDeployOptions `yaml:"charts"`
	internalCharts []HelmDeployOptions
}

func GetConfig() (Config, error) {
	config := Config{
		Name:        "dev",
		ServerCount: 1,
		AgentCount:  2,
		ServerPorts: []int{80, 443, 32500},
		Domain:      "localdev.me",
	}

	err := yaml.Unmarshal(chartByte, &config.internalCharts)
	if err != nil {
		return config, err
	}

	if f, err := os.Stat(confName); err == nil {
		config.load(f.Name())
	} else {
		hd, _ := os.UserHomeDir()

		if fi, err := os.Stat(filepath.Join(hd, confName)); err == nil {
			config.load(fi.Name())
		}
	}

	if len(config.ChartNames) > 0 {
		for _, chart := range config.internalCharts {
			config.Charts = append(config.Charts, chart)
		}
	} else {
		config.Charts = append(config.Charts, config.internalCharts...)
	}

	sort.SliceStable(config.Charts, func(i, j int) bool {
		return config.Charts[i].Stage < config.Charts[j].Stage
	})
	return config, nil
}

func (conf *Config) load(path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, conf)

}

func (config *Config) GetFlagSet() *pflag.FlagSet {
	flagset := pflag.NewFlagSet("clusterOptions", pflag.ExitOnError)
	flagset.StringVarP(&config.Name, "name", "n", config.Name, "name for the local cluster")
	flagset.StringVarP(&config.Domain, "domain", "d", config.Domain, "domain name for the local cluster")
	flagset.IntVarP(&config.ServerCount, "servers", "s", config.ServerCount, "number of servers for the cluster")
	flagset.IntVarP(&config.AgentCount, "agents", "a", config.AgentCount, "number of agents for the cluster")
	return flagset
}

func Create(streamOut bool, input Config) error {
	confPath, err := getAssetPath("assets/k3d/k3d-config.go.yaml", input)
	if err != nil {
		return err
	}
	addArgs := []string{
		"cluster", "create", input.Name,
		"--wait", "--config", confPath,
	}
	return utils.Exec(streamOut, utils.BinK3D, addArgs...)
}

func (config Config) GetChart(name string) *HelmDeployOptions {

	for _, ch := range config.internalCharts {
		if ch.Name == name {
			return &ch
		}
	}
	return nil
}
