package actor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPID(t *testing.T) {
	pid1, pid2 := GeneratePIDs("tester")
	assert.True(t, pid1.IsEqual(pid1))
	assert.False(t, pid2.IsEqual(pid1))
	assert.Contains(t, pid2.String(), ".tester")
}
