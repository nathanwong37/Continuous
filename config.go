package temp

import (
	"github.com/hashicorp/memberlist"
)

//MessengerConfig are configurations for messenger
type MessengerConfig struct {
	//memberlist config, has advertising port/ address and bind port / address
	memberConfig *memberlist.Config

	//Specific port that the RPC should listen on
	RPCPort int

	//For Local Connections only, can use different port for grpc
	LocalConnect bool
}

//Default for RPC
const (
	DefaultRPCPort int = 51284
)

//DefaultConfig returns a config pointer with default settings
func DefaultConfig() *MessengerConfig {
	return &MessengerConfig{
		memberConfig: memberlist.DefaultLocalConfig(),
		RPCPort:      DefaultRPCPort,
		LocalConnect: false,
	}
}

//DefaultWANConfig returns a config pointer with WAN default settings
func DefaultWANConfig() *MessengerConfig {
	return &MessengerConfig{
		memberConfig: memberlist.DefaultWANConfig(),
		RPCPort:      DefaultRPCPort,
		LocalConnect: false,
	}
}

//DefaultLANConfig returns a config pointer with default LAN settings
func DefaultLANConfig() *MessengerConfig {
	return &MessengerConfig{
		memberConfig: memberlist.DefaultLANConfig(),
		RPCPort:      DefaultRPCPort,
		LocalConnect: false,
	}
}

//CustomConfig is for any memberlist configs
func CustomConfig(config *memberlist.Config, isLocal bool) *MessengerConfig {
	if config == nil {
		config = memberlist.DefaultLocalConfig()
	}
	return &MessengerConfig{
		memberConfig: config,
		RPCPort:      DefaultRPCPort,
		LocalConnect: isLocal,
	}

}
