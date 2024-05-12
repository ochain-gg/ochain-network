package types

type OChainState struct {
	Id                 int    `cbor:"1,keyasint"`
	Size               int64  `cbor:"2,keyasint"`
	Height             int64  `cbor:"3,keyasint"`
	Hash               []byte `cbor:"4,keyasint"`
	LatestPortalUpdate uint64 `cbor:"5,keyasint"`
}

func (state *OChainState) SetHeight(height int64) {
	state.Height = height
}

func (state *OChainState) IncSize() {
	state.Size = state.Size + 1
}
