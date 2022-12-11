package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramstypes.ParamSet = (*Params)(nil)

var (
	KeyFeeRate = []byte("FeeRate")
	KeyMinInitialLiquidity = []byte("MinInitialLiquidity")
)

var (
	DefaultFeeRate = sdk.NewDecWithPrec(3, 3)
	DefaultMinInitialLiquidity = sdk.NewInt(1000)
)

func (params Params) String() string {
	out, _ := yaml.Marshal(params)
	return string(out)
}

func (params Params) Validate() error {
	err := validateFeeRate(params.FeeRate)
	if err != nil{
		return err
	}
	err = validateMinInitialLiquidity(params.MinInitialLiquidity)
	if err != nil{
		return err
	}
	return nil
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(feeRate sdk.Dec, minInitialLiquidity sdk.Int) Params {
	return Params{
		FeeRate: feeRate,
		MinInitialLiquidity: minInitialLiquidity,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultFeeRate, DefaultMinInitialLiquidity)
}

func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyFeeRate, &params.FeeRate, validateFeeRate),
		paramstypes.NewParamSetPair(KeyMinInitialLiquidity, &params.MinInitialLiquidity, validateMinInitialLiquidity),
	}
}

func validateFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("fee rate must not be negative: %s", v)
	}
	return nil
}

func validateMinInitialLiquidity(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("min initial liquidity must not be negative: %s", v)
	}
	return nil
}