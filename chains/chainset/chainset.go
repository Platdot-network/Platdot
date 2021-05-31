package chainset

const (
	NoneLike = iota
	EthLike
	PlatonLike

	SubLike
	PolkadotLike
	KusamaLike

	// ChainXV1Like ChainX V1 use Address Type
	ChainXV1Like
	ChainXAssetV1Like

	ChainXLike
	ChainXAssetLike
)

/// Chain name constants
const (
	NameUnimplemented		string = "unimplemented"

	NamePlaton				string = "platon"
	NameAlaya				string = "alaya"

	NameKusama				string = "kusama"
	NamePolkadot			string = "polkadot"
)

const(
	TokenATP	string = "ATP"
	TokenLAT	string = "LAT"

	TokenDOT	string = "DOT"
	TokenKSM	string = "KSM"
	TokenPCX	string = "PCX"

	TokenXBTC	string = "XBTC"
	TokenXAsset string = "XASSET"
)

type ChainInfo struct{
	Prefix 					string
	NativeToken 			string
	Type 					ChainType
}

var (
	ChainSets = [...]ChainInfo{
		{ NameAlaya, 			TokenATP,	PlatonLike},
		{ NamePlaton,			TokenLAT, PlatonLike},
		{ NameKusama, 		TokenKSM, KusamaLike },
		{ NamePolkadot,		TokenDOT, PolkadotLike },
	}
)
