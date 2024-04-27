package config

var (
	DefaultEVMRpc           = "http://localhost:8545/"
	DefaultEVMChainId       = uint64(31337)
	DefaultEVMPortalAddress = "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318"
)

type OChainConfig struct {
	EVMRpc           string
	EVMChainId       uint64
	EVMPortalAddress string
}

func DefaultConfig() *OChainConfig {
	return &OChainConfig{
		EVMRpc:           DefaultEVMRpc,
		EVMChainId:       DefaultEVMChainId,
		EVMPortalAddress: DefaultEVMPortalAddress,
	}
}
