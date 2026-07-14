package errgroupx

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func RunAll(ctx context.Context, n int, tasks []func(ctx context.Context) error) error {
	if n < 1 {
		panic("errgroupx.RunAll: limit must be positive")
	}
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(n)
	for _, task := range tasks {
		g.Go(func() error {
			return task(ctx)
		})
	}
	return g.Wait()
}
