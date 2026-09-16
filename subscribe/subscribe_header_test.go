package subscribe

import (
	"go-ethereum-homework/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscribeIntegration(t *testing.T) {
	ctx, client, _ := testutil.IntegrationWSSClient(t)
	err := Subscribe(ctx, client, 1)
	require.NoError(t, err)
}
