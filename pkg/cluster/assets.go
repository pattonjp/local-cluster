package cluster

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"text/template"
)

//go:embed assets/*
var assets embed.FS

func createTempFile(b []byte) (string, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	f.Write(b)
	f.Close()
	return f.Name(), nil
}

func getAssetPath(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return createTempFile(buf.Bytes())
}

func GetValuesFilesFor(chart string) ([]string, error) {
	files := []string{}
	prefix := fmt.Sprintf("assets/%s", chart)

	for _, tmp := range tpl.Templates() {
		if !strings.HasPrefix(tmp.Name(), prefix) {
			continue
		}

		isValFile := strings.HasSuffix(tmp.Name(), "values.yaml")
		isValFile = isValFile || strings.HasSuffix(tmp.Name(), "values.yml")
		isValFile = isValFile || strings.HasSuffix(tmp.Name(), "values.go.yaml")
		isValFile = isValFile || strings.HasSuffix(tmp.Name(), "values.go.yml")

		if isValFile {
			files = append(files, tmp.Name())
		}
	}
	return files, nil
}
func GetNonValuesFilesFor(chart string) ([]string, error) {
	files := []string{}
	prefix := fmt.Sprintf("assets/%s", chart)
	for _, tmp := range tpl.Templates() {
		if !strings.HasPrefix(tmp.Name(), prefix) {
			continue
		}
		notValues := !strings.HasSuffix(tmp.Name(), "values.yaml")
		notValues = notValues && !strings.HasSuffix(tmp.Name(), "values.yml")
		notValues = notValues && !strings.HasSuffix(tmp.Name(), "values.go.yaml")
		notValues = notValues && !strings.HasSuffix(tmp.Name(), "values.go.yml")

		if notValues {
			files = append(files, tmp.Name())
		}
	}
	return files, nil
}

var tpl *template.Template

func init() {
	tpl = template.New("")
	walker := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			t := tpl.New(path)
			buf, err := assets.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			_, e := t.Parse(string(buf))
			if e != nil {
				panic(e)
			}

		}
		return nil
	}
	fs.WalkDir(assets, ".", walker)
}
