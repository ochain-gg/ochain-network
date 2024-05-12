package types

type OChainValidator struct {
	Id                        uint64 `cbor:"1,keyasint"`
	Stacker                   string `cbor:"2,keyasint"`
	Validator                 string `cbor:"3,keyasint"`
	PublicKey                 string `cbor:"4,keyasint"`
	Enabled                   bool   `cbor:"5,keyasint"`
	StackingTreansactionHash  string `cbor:"6,keyasint"`
	UnstackingTransactionHash string `cbor:"7,keyasint"`
}
