package chainset

import (
	"fmt"

	"github.com/Platdot-Network/substrate-go/expand/chainx/xevents"
	"github.com/hacpy/go-ethereum/common"
	"github.com/rjman-ljm/platdot-utils/msg"
)

// The Eth-Like precision is 18 bits.
//var SingleEthLike = 1e18

// The currency of Sub-Like
const (
	SingleKSM int64 = 1e12
	SingleDOT int64 = 1e10
	SinglePCX int64 = 1e8
)

// The precision-difference constant between Eth-Like and Sub-Like
const (
	DiffKSM    = 1000000     /// KSM    is 12 digits
	DiffDOT    = 100000000   /// DOT    is 10 digits
	DiffXBTC   = 10000000000 /// XBTC   is 8  digits
	DiffPCX    = 10000000000 /// PCX	   is 8  digits
	DiffXAsset = 10000000000 /// XAsset is 8  digits
)

/// Fixed handling fee for cross-chain transactions
const (
	FixedKSMFee = SingleKSM * 1 / 100 /// 0.01KSM
	FixedDOTFee = 0 ///SingleDOT * 6 / 10  /// 0.6DOT
	FixedPCXFee = SinglePCX * 1 / 10  /// 0.1PCX
)

// ExtraFeeRate Additional formalities rate excluding fixed handling fee
const (
	ExtraFeeRate     int64 = 1000
	ExtraNoneFeeRate int64 = 0
)

type Currency struct {
	/// Set the token of the native token to zero
	AssetId      xevents.AssetId
	ResourceId   string
	Name         string
	Difference   int64
	FixedFee     int64
	ExtraFeeRate int64
}

var currencies = []Currency{
	{OriginAsset, ResourceIdOrigin, TokenKSM, DiffKSM, FixedKSMFee, ExtraFeeRate},
	{OriginAsset, ResourceIdOrigin, TokenDOT, DiffDOT, FixedDOTFee, ExtraFeeRate},
	{OriginAsset, ResourceIdOrigin, TokenPCX, DiffPCX, FixedPCXFee, ExtraFeeRate},
	{AssetXBTC, ResourceIdXBTC, TokenXBTC, DiffXBTC, 0, ExtraNoneFeeRate},
	{XAssetId, ResourceIdXAsset, TokenXAsset, DiffXAsset, 0, ExtraNoneFeeRate},
}

/// AssetId Type
const (
	OriginAsset xevents.AssetId = 0
	AssetXBTC   xevents.AssetId = 1
	XAssetId    xevents.AssetId = 999
)

const ResourceIdPrefix = "0000000000000000000000000000000000000000000000000000000000000"
const (
	ResourceIdOrigin string = ResourceIdPrefix + "000"
	ResourceIdXBTC   string = ResourceIdPrefix + "001"
	ResourceIdXAsset string = ResourceIdPrefix + "999"
	ResourceIdAKSM   string = ResourceIdPrefix + "000"
	ResourceIdPDOT   string = ResourceIdPrefix + "002"
)

func (bc *ChainCore) GetCurrencyByAssetId(assetId xevents.AssetId) (*Currency, error) {
	/// If token has assetId, return ChainX currency
	for _, currency := range currencies {
		if assetId != 0 && assetId == currency.AssetId {
			return &currency, nil
		} else if assetId == 0 && bc.ChainInfo.NativeToken == currency.Name {
			/// If token is native token, check the from chain
			return &currency, nil
		}
	}

	/// Currency not found
	return nil, fmt.Errorf("unimplemented currency")
}

func (bc *ChainCore) ConvertStringToResourceId(rId string) msg.ResourceId {
	return msg.ResourceIdFromSlice(common.FromHex(rId))
}

func (bc *ChainCore) GetCurrencyByResourceId(rId msg.ResourceId) (*Currency, error) {
	for _, currency := range currencies {
		if rId != bc.ConvertStringToResourceId(ResourceIdOrigin) && rId == bc.ConvertStringToResourceId(currency.ResourceId) {
			return &currency, nil
		} else if rId == bc.ConvertStringToResourceId(ResourceIdOrigin) && bc.ChainInfo.NativeToken == currency.Name {
			/// If token is native token, check the from chain
			return &currency, nil
		}
	}

	/// Currency not found
	return nil, fmt.Errorf("unimplemented currency")
}

func (bc *ChainCore) ConvertResourceIdToAssetId(rId msg.ResourceId) (xevents.AssetId, error) {
	for _, currency := range currencies {
		if rId != bc.ConvertStringToResourceId(ResourceIdOrigin) && rId == bc.ConvertStringToResourceId(currency.ResourceId) {
			return currency.AssetId, nil
		} else if rId == bc.ConvertStringToResourceId(ResourceIdOrigin) && bc.ChainInfo.NativeToken == currency.Name {
			/// If token is native token, check the from chain
			return currency.AssetId, nil
		}
	}

	/// Currency not found
	return xevents.AssetId(0), fmt.Errorf("unimplemented currency")
}
