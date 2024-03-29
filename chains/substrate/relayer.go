package substrate

import (
	"github.com/Platdot-Network/go-substrate-rpc-client/v3/signature"
	"github.com/Platdot-Network/go-substrate-rpc-client/v3/types"
)

type Relayer struct {
	kr                signature.KeyringPair
	otherSignatories  []types.AccountID
	totalRelayers     uint64
	multiSigThreshold uint16
	relayerId         uint64
	maxWeight         uint64
}

func NewRelayer(kr signature.KeyringPair, otherSignatories []types.AccountID, totalRelayers uint64,
	multiSigThreshold uint16, relayerId uint64, maxWeight uint64) Relayer {
	return Relayer{
		kr:                kr,
		otherSignatories:  otherSignatories,
		totalRelayers:     totalRelayers,
		multiSigThreshold: multiSigThreshold,
		relayerId:         relayerId,
		maxWeight:         maxWeight,
	}
}
