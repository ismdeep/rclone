package bvfs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_MkDir(t *testing.T) {
	assert.NoError(t, NewClient().SignIn().MkDir(context.Background(), "Wallpapers/Linux"))
}
