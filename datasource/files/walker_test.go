package files

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestExtract(t *testing.T) {
	t.Skip("only local testing")

	c := make(chan File)

	wg, ctx := errgroup.WithContext(context.TODO())

	wg.Go(func() error {
		return Extract(ctx, "/Users/nvrusin/black", c)
	})

	wg.Go(func() error {
		for f := range c {
			t.Log("file", f.Package, f.Name, f.Lines, f.Symbols)
		}

		return nil
	})

	require.NoError(t, wg.Wait())
}
