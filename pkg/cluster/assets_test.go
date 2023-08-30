package cluster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	// getAssetPath(name string)
	l, err := assets.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range l {
		fmt.Println("dir:", v.Name())
	}

}

func TestListAllFiles(t *testing.T) {
	files, err := GetNonValuesFilesFor("traefik")
	assert.NoError(t, err)
	assert.Len(t, files, 2)

	for _, tmp := range tpl.Templates() {
		fmt.Println(tmp.Name())
	}
}

func TestAllFiles(t *testing.T) {
	cases := []struct {
		name string
	}{
		{"assets/kube-prometheus-stack/values.yaml"},
		{"assets/traefik/values.yaml"},
		{"assets/k3d/k3d-config.go.yaml"},
	}

	for _, tc := range cases {
		f, err := getAssetPath(tc.name, nil)

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(f)

	}

}

func TestK3DConfigFile(t *testing.T) {
	f, err := getAssetPath("assets/k3d/k3d-config.go.yaml", nil)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(f)

}
