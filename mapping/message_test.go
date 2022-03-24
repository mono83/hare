package mapping

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_prepareRoutingKey(t *testing.T) {
	m := Message{
		Headers: map[string]interface{}{
			"some":    455,
			"another": "bar",
		},
	}

	assert.Equal(t, "foo", m.prepareRoutingKey("foo"))
	assert.Equal(t, "$", m.prepareRoutingKey("$"))
	assert.Equal(t, "#", m.prepareRoutingKey("#"))
	assert.Equal(t, "@", m.prepareRoutingKey("@"))
	assert.Equal(t, "455", m.prepareRoutingKey("$some"))
	assert.Equal(t, "455", m.prepareRoutingKey("@some"))
	assert.Equal(t, "455", m.prepareRoutingKey("#some"))
	assert.Equal(t, "some", m.prepareRoutingKey("some"))
	assert.Equal(t, "bar", m.prepareRoutingKey("$another"))
	assert.Equal(t, "bar", m.prepareRoutingKey("@another"))
	assert.Equal(t, "bar", m.prepareRoutingKey("#another"))
	assert.Equal(t, "another", m.prepareRoutingKey("another"))

	assert.Equal(t, "", Message{}.prepareRoutingKey("$some"))
}
