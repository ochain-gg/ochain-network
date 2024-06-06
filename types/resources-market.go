package types

type OChainResourcesMarket struct {
	UniverseId string `cbor:"1,keyasint"`
	FeesRate   uint64 `cbor:"2,keyasint"`

	MetalReserveRatio uint64 `cbor:"3,keyasint"`
	MetalPoolBalance  uint64 `cbor:"4,keyasint"`
	MetalSupplyMinted uint64 `cbor:"5,keyasint"`

	CrystalReserveRatio uint64 `cbor:"6,keyasint"`
	CrystalPoolBalance  uint64 `cbor:"7,keyasint"`
	CrystalSupplyMinted uint64 `cbor:"8,keyasint"`

	DeuteriumReserveRatio uint64 `cbor:"9,keyasint"`
	DeuteriumPoolBalance  uint64 `cbor:"10,keyasint"`
	DeuteriumSupplyMinted uint64 `cbor:"11,keyasint"`
}
