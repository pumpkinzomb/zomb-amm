package types

func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultParams(), 0 , nil)
}

func NewGenesisState(params Params, lastPairId uint64, pairs []Pair) *GenesisState {
	return &GenesisState{
		Params: params,
		LastPairId: lastPairId,
		Pairs: pairs,
	}
}

func (gs GenesisState) Validate() error {
	err := gs.Params.Validate()
	if err != nil{
		return err
	}
	return nil
}