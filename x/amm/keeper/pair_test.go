package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	"github.com/pumpkinzomb/zomb-amm/utils"
	"github.com/pumpkinzomb/zomb-amm/x/amm/types"
)

func (s *KeeperTestSuite) TestAddAndRemoveLiquidity() {
	fromAddr := utils.SampleAddress(1)
	coins := sdk.NewCoins(
		sdk.NewInt64Coin("denom1", 1000000),
		sdk.NewInt64Coin("denom2", 100000000))
	s.Require().NoError(testutil.FundAccount(s.app.BankKeeper, s.ctx, fromAddr, coins))
	mintedShare, err := s.keeper.AddLiquidity(s.ctx, fromAddr, coins)
	s.Require().NoError(err)
	s.Require().Equal("10000000"+types.ShareDenomPrefix+"1", mintedShare.String())

	withdrawnCoins, err := s.keeper.RemoveLiquidity(
		s.ctx, fromAddr, sdk.NewInt64Coin(mintedShare.Denom, 5000000))
	s.Require().NoError(err)
	s.Require().Equal("500000denom1,50000000denom2", withdrawnCoins.String())
}

func (s *KeeperTestSuite) TestSwapExactIn() {
	pairCreatorAddr := utils.SampleAddress(1)
	coins := sdk.NewCoins(
		sdk.NewInt64Coin("denom1", 1000000),
		sdk.NewInt64Coin("denom2", 100000000))
	s.Require().NoError(testutil.FundAccount(s.app.BankKeeper, s.ctx, pairCreatorAddr, coins))
	_, err := s.keeper.AddLiquidity(s.ctx, pairCreatorAddr, coins)
	s.Require().NoError(err)

	fromAddr := utils.SampleAddress(2)
	s.Require().NoError(
		testutil.FundAccount(s.app.BankKeeper, s.ctx, fromAddr,
			sdk.NewCoins(sdk.NewInt64Coin("denom1", 1000000))))
	coinOut, err := s.keeper.SwapExactIn(
		s.ctx, fromAddr, sdk.NewInt64Coin("denom1", 10000), sdk.NewInt64Coin("denom2", 950000))
	s.Require().NoError(err)
	s.Require().Equal("987158denom2", coinOut.String())

	coinOut, err = s.keeper.SwapExactIn(
		s.ctx, fromAddr, sdk.NewInt64Coin("denom2", 987158), sdk.NewInt64Coin("denom1", 9000))
	s.Require().NoError(err)
	s.Require().Equal("9940denom1", coinOut.String())
}

func (s *KeeperTestSuite) TestSwapExactOut() {
	pairCreatorAddr := utils.SampleAddress(1)
	coins := sdk.NewCoins(
		sdk.NewInt64Coin("denom1", 1000000),
		sdk.NewInt64Coin("denom2", 100000000))
	s.Require().NoError(testutil.FundAccount(s.app.BankKeeper, s.ctx, pairCreatorAddr, coins))
	_, err := s.keeper.AddLiquidity(s.ctx, pairCreatorAddr, coins)
	s.Require().NoError(err)

	fromAddr := utils.SampleAddress(2)
	s.Require().NoError(
		testutil.FundAccount(s.app.BankKeeper, s.ctx, fromAddr,
			sdk.NewCoins(sdk.NewInt64Coin("denom1", 1000000))))
	coinIn, err := s.keeper.SwapExactOut(
		s.ctx, fromAddr, sdk.NewInt64Coin("denom2", 987158), sdk.NewInt64Coin("denom1", 12000))
	s.Require().NoError(err)
	s.Require().Equal("10000denom1", coinIn.String())

	coinIn, err = s.keeper.SwapExactOut(
		s.ctx, fromAddr, sdk.NewInt64Coin("denom1", 9940), sdk.NewInt64Coin("denom2", 987158))
	s.Require().NoError(err)
	s.Require().Equal("987081denom2", coinIn.String())
}
