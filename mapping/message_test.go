package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestStripUnsupportedPublishingHeaders(t *testing.T) {
	assert.Nil(t, stripUnsupportedPublishingHeaders(nil))
	assert.NotNil(t, stripUnsupportedPublishingHeaders(map[string]interface{}{}))

	datum := map[string]interface{}{
		"foo":     "bar",
		"baz":     1,
		"x-death": "xxx",
	}

	datum = stripUnsupportedPublishingHeaders(datum)
	assert.Len(t, datum, 2)
	assert.Equal(t, "bar", datum["foo"])
	assert.Equal(t, 1, datum["baz"])
}
