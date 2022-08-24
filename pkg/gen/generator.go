package gen

import (
	_ "embed"
	"io"
	"text/template"
)

type Condition struct {
	Regex  string
	Action string
}

var (
	//go:embed pac.tmpl
	templStr string

	templ *template.Template
)

func init() {
	templ = template.Must(template.New("pac").Parse(templStr))
}

func Generate(wr io.Writer, data []Condition) error {
	err := templ.Execute(wr, &data)
	return err
}
