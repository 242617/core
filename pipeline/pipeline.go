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

/*
New creates pipeline that call functions in this order:
  - Then...
  - ThenCatch
  - Else...
  - ElseCatch
  - Catch...
  - Invoke
  - Call

Example:

	errCh := make(chan error)
	go pipeline.New(ctx).
		Then(func(ctx context.Context) error {
			return nil
		}).
		Else(func(ctx context.Context) error {
			return nil
		}).
		Run(func(err error) { errCh <- err })
	return <-errCh
*/
func New(ctx context.Context, funcs ...Func) *Pipeline {
	return NewWithOptions(WithContext(ctx)).Then(funcs...)
}

type (
	Func       = func(context.Context) error
	CatchFunc  = func(error) error
	ErrFunc    = func(error)
	InvokeFunc = func()
	Pipeline   struct {
		ctx    context.Context
		err    error
		layers []layer
	}
	layer struct {
		funcs, fallbacks         []Func
		thenCatcher, elseCatcher CatchFunc
		catchers                 []CatchFunc
		before, after            InvokeFunc
		reset                    bool
	}
)

func (p *Pipeline) Then(funcs ...Func) *Pipeline {
	p.layers = append(p.layers, layer{funcs: funcs})
	return p
}

func (p *Pipeline) ThenCatch(f CatchFunc) *Pipeline {
	p.layers[len(p.layers)-1].thenCatcher = f
	return p
}

func (p *Pipeline) Else(fallbacks ...Func) *Pipeline {
	if p.layers[len(p.layers)-1].fallbacks == nil {
		p.layers[len(p.layers)-1].fallbacks = fallbacks
	}
	return p
}

func (p *Pipeline) ElseCatch(catcher CatchFunc) *Pipeline {
	p.layers[len(p.layers)-1].elseCatcher = catcher
	return p
}

func (p *Pipeline) Catch(catchers ...CatchFunc) *Pipeline {
	if p.layers[len(p.layers)-1].catchers == nil {
		p.layers[len(p.layers)-1].catchers = catchers
	}
	return p
}

func (p *Pipeline) Before(before InvokeFunc) *Pipeline {
	p.layers[len(p.layers)-1].before = before
	return p
}

func (p *Pipeline) After(after InvokeFunc) *Pipeline {
	p.layers[len(p.layers)-1].after = after
	return p
}

func (p *Pipeline) Reset() *Pipeline {
	p.layers = append(p.layers, layer{reset: true})
	return p
}

func (p *Pipeline) Run(errFunc ErrFunc) {
	for _, layer := range p.layers {
		if layer.reset {
			p.err = nil
			continue
		}
		if layer.before != nil {
			layer.before()
		}
		if p.err == nil && len(layer.funcs) > 0 {
			p.err = p.process(layer.funcs...)
			if p.err != nil && layer.thenCatcher != nil {
				p.err = p.intercept(layer.thenCatcher)
			}
			if p.err != nil && len(layer.fallbacks) > 0 {
				p.err = p.process(layer.fallbacks...)
			}
			if p.err != nil && layer.elseCatcher != nil {
				p.err = p.intercept(layer.elseCatcher)
			}
		}
		if layer.after != nil {
			layer.after()
		}
	}
	errFunc(p.err)
}

func (p *Pipeline) intercept(interceptor CatchFunc) error { return interceptor(p.err) }

func (p *Pipeline) process(funcs ...Func) error {
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
