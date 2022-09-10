package pipeline

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/sync/errgroup"
)

func NewWithOptions(options ...option) *Pipeline {
	var p Pipeline
	for _, option := range options {
		option(&p)
	}
	return &p
}

func New(ctx context.Context, funcs ...PipelineFunc) *Pipeline {
	return NewWithOptions(WithContext(ctx)).Then(funcs...)
}

type (
	PipelineFunc    = func(context.Context) error
	PipelineErrFunc = func(error)
	Pipeline        struct {
		ctx    context.Context
		err    error
		layers []layer
	}
	layer struct {
		funcs, fallbacks []PipelineFunc
		reset            bool
	}
)

func (p *Pipeline) Then(funcs ...PipelineFunc) *Pipeline {
	p.layers = append(p.layers, layer{funcs: funcs})
	return p
}

func (p *Pipeline) Else(fallbacks ...PipelineFunc) *Pipeline {
	if p.layers[len(p.layers)-1].fallbacks == nil {
		p.layers[len(p.layers)-1].fallbacks = fallbacks
	}
	return p
}

func (p *Pipeline) Reset() *Pipeline {
	p.layers = append(p.layers, layer{reset: true})
	return p
}

func (p *Pipeline) Run(errFunc PipelineErrFunc) {
	for _, layer := range p.layers {
		if layer.reset {
			p.err = nil
			continue
		}
		if p.err == nil && len(layer.funcs) > 0 {
			p.err = p.process(layer.funcs...)
			if p.err != nil && len(layer.fallbacks) > 0 {
				p.err = p.process(layer.fallbacks...)
			}
		}
	}
	errFunc(p.err)
}

func (p *Pipeline) process(funcs ...PipelineFunc) error {
	errCh := make(chan error)
	go func() {
		group, ctx := errgroup.WithContext(p.ctx)
		for _, f := range funcs {
			f := f
			group.Go(func() error { return f(ctx) })
		}
		errCh <- group.Wait()
		close(errCh)
	}()

	var err error
	select {
	case <-p.ctx.Done():
		err = p.ctx.Err()
	case err = <-errCh:
	}
	return err
}

func (p *Pipeline) Call(f func(...interface{}), args ...interface{}) *Pipeline {
	f(args...)
	return p
}
func (p *Pipeline) Invoke(f func()) *Pipeline {
	f()
	return p
}

func (p *Pipeline) String() string {
	var info strings.Builder
	info.WriteString("Pipeline: {\n")
	for i, layer := range p.layers {
		var layerInfo string
		if layer.reset {
			layerInfo = "reset"
		} else {
			layerInfo = fmt.Sprintf("funcs: %d, fallbacks: %d", len(layer.funcs), len(layer.fallbacks))
		}
		info.WriteString(fmt.Sprintf("[%2d]: %s\n", i, layerInfo))
	}
	if p.err != nil {
		info.WriteString(fmt.Sprintf("error: %q\n", p.err.Error()))
	}
	info.WriteString("}")
	return info.String()
}
