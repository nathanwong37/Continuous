package continuous

import (
	"testing"

	conf "github.com/Continuous/config"
	"github.com/stretchr/testify/require"
)

func TestGrpc(t *testing.T) {
	config := conf.DefaultConfig()
	test := NewMessenger(config)
	test.Join(nil)
	uuidstr, err := test.client.CreateTimer(70, "Nathan Wong", "00:00:10", "")
	require.NoError(t, err)
	_, err = test.client.DeleteTimer(uuidstr, "Nathan Wong")
	require.NoError(t, err)
}
