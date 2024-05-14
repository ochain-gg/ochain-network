package types

type OChainState struct {
	Size               int64  `cbor:"1,keyasint"`
	Height             int64  `cbor:"2,keyasint"`
	Hash               []byte `cbor:"3,keyasint"`
	LatestPortalUpdate uint64 `cbor:"4,keyasint"`
}

func (state *OChainState) SetHeight(height int64) {
	state.Height = height
}

func (state *OChainState) IncSize() {
	state.Size = state.Size + 1
}

func (state *OChainState) SetLatestPortalUpdate(block uint64) {
	state.LatestPortalUpdate = block
}
