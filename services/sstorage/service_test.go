package sstorage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceStorage(t *testing.T) {
	service := NewServiceStorage()
	require.NotNil(t, service)
}
