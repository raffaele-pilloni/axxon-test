package executor

import "context"

type Name string

type Interface interface {
	GetName() Name
	Run(ctx context.Context, args []string) error
}
