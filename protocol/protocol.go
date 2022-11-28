package protocol

import "context"

type Lifecycle interface {
	Start(context.Context) error
	Stop(context.Context) error
}
