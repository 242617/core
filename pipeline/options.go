package pipeline

import "context"

type option func(p *Pipeline)

func WithContext(ctx context.Context) option { return func(p *Pipeline) { p.ctx = ctx } }
