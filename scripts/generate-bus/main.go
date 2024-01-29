// main generates bus based on all structs in a file.
package main

import (
	"flag"
	"os"
	"path"
	"regexp"
	"text/template"
)

const rawTemplate = `// Code generated by {{ .By }}; DO NOT EDIT.
package {{ .Package }}

import (
{{- range .Imports }}
	"{{.}}"
{{- end }}
)

func busLogError(err error) {
	if err != nil {
		log.Err(err).Str("package", "{{ .Package }}").Msg("Failed to handle event")
	}
}

func NewBus() *Bus {
	return &Bus{
		ServiceContext: sutureext.NewServiceContext("{{ .Package }}.Bus"),
	}
}

type Bus struct {
	sutureext.ServiceContext
{{- range .Events }}
	on{{.}} []func(ctx context.Context, event {{ $.EventPackage }}{{ . }}) error
{{- end }}
}

func (b *Bus) Register(pub pubsub.Pub) (*Bus) {
{{- range .Events }}
	b.On{{ . }}(func(ctx context.Context, evt {{ $.EventPackage }}{{ . }}) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
{{- end }}
	return b
}

{{ range .Events }}
func (b *Bus) On{{ . }}(h func(ctx context.Context, evt {{ $.EventPackage }}{{ . }}) error) {
	b.on{{ . }} = append(b.on{{ . }}, h)
}
{{ end }}

{{ range .Events }}
func (b *Bus) {{ . }}(evt {{ $.EventPackage }}{{ . }}) {
	for _, v := range b.on{{ . }} {
		busLogError(v(b.Context(), evt))
	}
}
{{ end }}
`

type Data struct {
	By           string
	Package      string
	Imports      []string
	EventPackage string
	Events       []string
}

func main() {
	outputFilePath := "./internal/event/bus.gen.go"
	inputFilePath := "./internal/event/event.go"

	flag.StringVar(&outputFilePath, "output", "", "")
	flag.StringVar(&inputFilePath, "input", "", "")

	flag.Parse()

	outputFilePath = path.Clean(outputFilePath)
	inputFilePath = path.Clean(inputFilePath)

	_ = os.Remove(outputFilePath)

	var events []string
	for _, v := range must2(regexp.Compile(`type (.*?) struct {`)).FindAllStringSubmatch(string(must2(os.ReadFile(inputFilePath))), -1) {
		events = append(events, v[1])
	}
	data := Data{
		By:      "generate-bus.go",
		Package: "event",
		Imports: []string{
			"context",
			"errors",
			"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub",
			"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext",
			"github.com/rs/zerolog/log",
		},
		EventPackage: "",
		Events:       events,
	}

	templ := must2(template.New("").Parse(rawTemplate))

	file := must2(os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY, 0644))

	must(templ.Execute(file, data))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2[T any](d T, err error) T {
	if err != nil {
		panic(err)
	}
	return d
}
