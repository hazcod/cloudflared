package pogs

import (
	"reflect"
	"testing"

	"github.com/cloudflare/cloudflared/tunnelrpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	capnp "zombiezen.com/go/capnproto2"
)

// Assert *SystemName implements Scope
var _ Scope = (*SystemName)(nil)

// Assert *Group implements Scope
var _ Scope = (*Group)(nil)

func TestScope(t *testing.T) {
	testCases := []Scope{
		&SystemName{systemName: "my_system"},
		&Group{group: "my_group"},
	}
	for i, testCase := range testCases {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		capnpEntity, err := tunnelrpc.NewScope(seg)
		if !assert.NoError(t, err) {
			t.Fatal("Couldn't initialize a new message")
		}
		err = MarshalScope(capnpEntity, testCase)
		if !assert.NoError(t, err, "testCase index %v failed to marshal", i) {
			continue
		}
		result, err := UnmarshalScope(capnpEntity)
		if !assert.NoError(t, err, "testCase index %v failed to unmarshal", i) {
			continue
		}
		assert.Equal(t, testCase, result, "testCase index %v didn't preserve struct through marshalling and unmarshalling", i)
	}
}

func TestConnectParameters(t *testing.T) {
	testCases := []*ConnectParameters{
		sampleConnectParameters(),
		sampleConnectParameters(func(c *ConnectParameters) {
			c.Scope = &SystemName{systemName: "my_system"}
		}),
		sampleConnectParameters(func(c *ConnectParameters) {
			c.Tags = nil
		}),
	}
	for i, testCase := range testCases {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		capnpEntity, err := tunnelrpc.NewCapnpConnectParameters(seg)
		if !assert.NoError(t, err) {
			t.Fatal("Couldn't initialize a new message")
		}
		err = MarshalConnectParameters(capnpEntity, testCase)
		if !assert.NoError(t, err, "testCase index %v failed to marshal", i) {
			continue
		}
		result, err := UnmarshalConnectParameters(capnpEntity)
		if !assert.NoError(t, err, "testCase index %v failed to unmarshal", i) {
			continue
		}
		assert.Equal(t, testCase, result, "testCase index %v didn't preserve struct through marshalling and unmarshalling", i)
	}
}

func sampleConnectParameters(overrides ...func(*ConnectParameters)) *ConnectParameters {
	cloudflaredID, err := uuid.Parse("ED7BA470-8E54-465E-825C-99712043E01C")
	if err != nil {
		panic(err)
	}
	sample := &ConnectParameters{
		OriginCert:          []byte("my-origin-cert"),
		CloudflaredID:       cloudflaredID,
		NumPreviousAttempts: 19,
		Tags: []Tag{
			Tag{
				Name:  "provision-method",
				Value: "new",
			},
		},
		CloudflaredVersion: "7.0",
		Scope:              &Group{group: "my_group"},
	}
	sample.ensureNoZeroFields()
	for _, f := range overrides {
		f(sample)
	}
	return sample
}

func (c *ConnectParameters) ensureNoZeroFields() {
	ensureNoZeroFieldsInSample(reflect.ValueOf(c), []string{})
}
