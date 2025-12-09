package diagram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/nats-io/nats.go/micro"
	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2oracle"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

type GraphOpts struct {
	GenerateSVG bool
	EdgeOpts    []EdgeOpt
}
type EdgeOpt func(*d2graph.Graph, string) (*d2graph.Graph, error)

func WithAnimation(graph *d2graph.Graph, edge string) (*d2graph.Graph, error) {
	a := fmt.Sprintf("%s.style.animated", edge)
	return d2oracle.Set(graph, nil, a, nil, go2.Pointer("true"))
}

func CreateEdge(graph *d2graph.Graph, to, from string, opts ...EdgeOpt) (*d2graph.Graph, error) {
	edge := fmt.Sprintf("%s -> %s", to, from)
	graph, key, _ := d2oracle.Create(graph, nil, edge)
	var err error
	for _, v := range opts {
		graph, err = v(graph, key)
		if err != nil {
			return nil, err
		}
	}

	return graph, nil
}

func New(ctx context.Context, logger *slog.Logger, micro micro.Info, opts GraphOpts) (string, []byte, error) {
	dCtx := log.With(context.Background(), logger)
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2elklayout.DefaultLayout, nil
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}
	renderOpts := &d2svg.RenderOpts{
		Pad: go2.Pointer(int64(5)),
	}
	_, graph, _ := d2lib.Compile(dCtx, "NATS", compileOpts, nil)

	graph, m, _ := d2oracle.Create(graph, nil, micro.Name)

	for _, v := range micro.Endpoints {
		var err error

		sub := v.Subject
		name := v.Name
		if strings.Contains(v.Subject, ".") {
			sub = fmt.Sprintf(`"%s"`, v.Subject)
		}
		graph, _, _ = d2oracle.Create(graph, nil, name)
		graph, _, _ = d2oracle.Create(graph, nil, sub)
		graph, err = CreateEdge(graph, name, sub, opts.EdgeOpts...)
		if err != nil {
			return "", nil, err
		}
		graph, _ = d2oracle.Move(graph, nil, name, fmt.Sprintf("%s.%s", m, name), false)
		graph, _ = d2oracle.Move(graph, nil, sub, fmt.Sprintf("NATS.%s", sub), false)
	}

	diagramStr := d2format.Format(graph.AST)
	diagram, _, _ := d2lib.Compile(dCtx, d2format.Format(graph.AST), compileOpts, renderOpts)

	if opts.GenerateSVG {
		out, err := d2svg.Render(diagram, renderOpts)
		if err != nil {
			return "", nil, err
		}

		return "", out, nil
	}

	return diagramStr, nil, nil
}
