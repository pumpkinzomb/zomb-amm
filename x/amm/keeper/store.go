package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pumpkinzomb/zomb-amm/x/amm/types"
)

func (k Keeper) GetLastPairID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastPairIDKey)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastPairID(ctx sdk.Context, lastPairID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastPairIDKey, sdk.Uint64ToBigEndian(lastPairID))
}

func (k Keeper) GetPair(ctx sdk.Context, pairID uint64) (pair types.Pair, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPairKey(pairID))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshal(bz, &pair)
	return pair, true
}

func (k Keeper) SetPair(ctx sdk.Context, pair types.Pair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&pair)
	store.Set(types.GetPairKey(pair.Id), bz)
}

func (k Keeper) DeletePair(ctx sdk.Context, pair types.Pair) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPairKey(pair.Id))
}

func (k Keeper) IterateAllPairs(ctx sdk.Context, cb func(pair types.Pair) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PairKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var pair types.Pair
		k.cdc.MustUnmarshal(iter.Value(), &pair)

		if cb(pair) {
			break
		}
	}
}

func (k Keeper) GetAllPairs(ctx sdk.Context) (pairs []types.Pair) {
	k.IterateAllPairs(ctx, func(pair types.Pair) (stop bool){
		pairs = append(pairs, pair)
		return false
	})
	return
}

func (k Keeper) GetPairByDenoms(ctx sdk.Context, denomA, denomB string) (pair types.Pair, found bool) {
	denom0, denom1 := types.SortDenoms(denomA, denomB)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPairIndexKey(denom0, denom1))
	if bz == nil {
		return 
	}
	return k.GetPair(ctx, sdk.BigEndianToUint64(bz))
}

func (k Keeper) SetPairIndex(ctx sdk.Context, pair types.Pair) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(pair.Id)
	store.Set(types.GetPairIndexKey(pair.Denom0, pair.Denom1), bz)
}

func (k Keeper) DeletePairIndex(ctx sdk.Context, pair types.Pair) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPairIndexKey(pair.Denom0, pair.Denom1))
}