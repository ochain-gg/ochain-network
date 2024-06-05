package types

type OChainValidator struct {
	Id                        uint64 `cbor:"1,keyasint"`
	Stacker                   string `cbor:"2,keyasint"` //0x address of portal wich stack OCT Tokens
	Validator                 string `cbor:"3,keyasint"` //0x address of the validator
	PublicKey                 string `cbor:"4,keyasint"` //ed public key of the validator
	Enabled                   bool   `cbor:"5,keyasint"`
	StackingTransactionHash   string `cbor:"6,keyasint"`
	UnstackingTransactionHash string `cbor:"7,keyasint"`
}
