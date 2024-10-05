package pipeline

import "context"

type option func(p *Pipeline)

func WithContext(ctx context.Context) option { return func(p *Pipeline) { p.ctx = ctx } }

func withError(err error) option {
	return func(p *Pipeline) { p.err = err }
}

func withLayers(layers ...layer) option {
	return func(p *Pipeline) { p.layers = append(p.layers, layers...) }
}
