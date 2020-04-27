package build

import (
	"bytes"
	"fmt"
	"text/template"
)

var (
	// Version holds the tag of the build
	Version = ""
	// Commit holds the git commit sha
	Commit = ""
	// Date is the build date
	Date = ""
)

const versionTpl = `
Version:    v{{ .Version }}
SHA:        {{ .Commit }}
Built On:   {{ .Date }}`

// PrintVersion outputs the short version info
func PrintVersion() {
	if Version != "" {
		fmt.Printf("v%s\n", Version)
	}
}

// PrintLongVersion outputs the full version info
func PrintLongVersion() error {
	if Version == "" {
		return nil
	}

	data := struct {
		Version string
		Commit  string
		Date    string
	}{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}

	var tpl bytes.Buffer

	t, err := template.New("build").Parse(versionTpl)
	if err != nil {
		return err
	}

	if err := t.Execute(&tpl, data); err != nil {
		return err
	}

	fmt.Println(tpl.String())
	return nil
}
