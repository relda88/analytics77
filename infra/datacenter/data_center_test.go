package datacenter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataCenter(t *testing.T) {
	dc := NewDataCenter()
	require.NotNil(t, dc)

	require.Empty(t, dc.AddEvents())
	require.Empty(t, dc.AddEvents(nil))

}
