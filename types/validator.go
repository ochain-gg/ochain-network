package types

type OChainValidatorWithAttributes struct {
	Id                        uint64 `cbor:"id"`
	Stacker                   string `cbor:"stacker"`   //0x address of portal wich stack OCT Tokens
	Validator                 string `cbor:"validator"` //0x address of the validator
	PublicKey                 string `cbor:"publicKey"` //ed public key of the validator
	Power                     int64  `cbor:"power"`
	Enabled                   bool   `cbor:"enabled"`
	StackingTransactionHash   string `cbor:"stackingTransactionHash"`
	UnstackingTransactionHash string `cbor:"unstackingTransactionHash"`
}

type OChainValidator struct {
	Id                        uint64 `cbor:"1,keyasint"`
	Stacker                   string `cbor:"2,keyasint"` //0x address of portal wich stack OCT Tokens
	Validator                 string `cbor:"3,keyasint"` //0x address of the validator
	PublicKey                 string `cbor:"4,keyasint"` //ed public key of the validator
	Power                     int64  `cbor:"5,keyasint"`
	Enabled                   bool   `cbor:"6,keyasint"`
	StackingTransactionHash   string `cbor:"7,keyasint"`
	UnstackingTransactionHash string `cbor:"8,keyasint"`
}

func (val *OChainValidator) WithAttributes() OChainValidatorWithAttributes {
	return OChainValidatorWithAttributes{
		Id:                        val.Id,
		Stacker:                   val.Stacker,
		Validator:                 val.Validator,
		PublicKey:                 val.PublicKey,
		Power:                     val.Power,
		Enabled:                   val.Enabled,
		StackingTransactionHash:   val.StackingTransactionHash,
		UnstackingTransactionHash: val.UnstackingTransactionHash,
	}
}
