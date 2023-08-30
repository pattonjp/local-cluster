package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChartsOut(t *testing.T) {
	config, err := GetConfig()
	assert.NoError(t, err, "getting config")

	assert.GreaterOrEqual(t, len(config.Charts), 1, "charts found")

}

func TestChartsRead(t *testing.T) {
	config, err := GetConfig()
	assert.NoError(t, err, "getting config")
	ch := config.GetChart("traefik")

	assert.NotNil(t, ch)
	assert.Equal(t, ch.AddCerts, true)
	assert.Equal(t, base, ch.Stage, "stage")
	assert.Equal(t, "operations", ch.Namespace)
	assert.Equal(t, "traefik/traefik", ch.Chart)
	assert.Equal(t, "traefik", ch.RepoName)
	assert.Equal(t, "https://traefik.github.io/charts", ch.RepoURL)

}

func TestChartValuesFiles(t *testing.T) {
	f, err := GetValuesFilesFor("loki")
	assert.NoError(t, err, "error finding loki values file")
	assert.Len(t, f, 1, "loki values file")
}
