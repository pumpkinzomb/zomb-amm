package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chaosapp "github.com/pumpkinzomb/zomb-amm/app"
	"github.com/pumpkinzomb/zomb-amm/x/amm/keeper"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	app    *chaosapp.App
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chaosapp.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.keeper = s.app.AMMKeeper
}
