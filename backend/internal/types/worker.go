package types

import "context"

type Job interface {
	Execute(ctx context.Context) error
	ID() string
}

type Submitter interface {
	Submit(job Job)
}
