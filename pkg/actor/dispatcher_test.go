package actor_test

import (
	"testing"

	"github.com/pix303/actor-lib/pkg/actor"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	// pid1, pid2 := GeneratePIDs()
	a := GenerateActor("test")
	b := GenerateActor("test1")
	d := actor.NewActorDispatcher()
	_ = d.RegisterActor(a)
	numActors := d.RegisterActor(b)
	assert.Equal(t, numActors, 2)
}
