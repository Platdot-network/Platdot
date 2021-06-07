// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package ethlike

import (
	"github.com/ChainSafe/log15"
	"github.com/Platdot-network/Platdot/bindings/Bridge"
	"github.com/Platdot-network/Platdot/chains/chainset"
	"github.com/rjman-ljm/platdot-utils/core"
	"github.com/rjman-ljm/platdot-utils/crypto/secp256k1"
	metrics "github.com/rjman-ljm/platdot-utils/metrics/types"
	"github.com/rjman-ljm/platdot-utils/msg"
)

var _ core.Writer = &writer{}

var PassedStatus uint8 = 2
var TransferredStatus uint8 = 3
var CancelledStatus uint8 = 4

type writer struct {
	cfg            Config
	conn           Connection
	bridgeContract *Bridge.Bridge // instance of bound receiver bridgeContract
	kp             secp256k1.Keypair
	log            log15.Logger
	stop           <-chan int
	sysErr         chan<- error // Reports fatal error to core
	metrics        *metrics.ChainMetrics
	chainCore      *chainset.ChainCore
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *Config, log log15.Logger, kp secp256k1.Keypair, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics, bc *chainset.ChainCore) *writer {
	return &writer{
		cfg:       *cfg,
		conn:      conn,
		kp:        kp,
		log:       log,
		stop:      stop,
		sysErr:    sysErr,
		metrics:   m,
		chainCore: bc,
	}
}

func (w *writer) start() error {
	w.log.Debug("Starting writer...", "chain", w.cfg.name)
	return nil
}

// setContract adds the bound receiver bridgeContract to the writer
func (w *writer) setContract(bridge *Bridge.Bridge) {
	w.bridgeContract = bridge
}

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success, this should be ignored except for within tests.
func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "recipient", m.Payload[1])
	switch m.Type {
	case msg.MultiSigTransfer:
		return w.createMultiSigProposal(m)
	case msg.FungibleTransfer:
		return w.createErc20Proposal(m)
	case msg.NonFungibleTransfer:
		return w.createErc721Proposal(m)
	case msg.GenericTransfer:
		return w.createGenericDepositProposal(m)
	default:
		w.log.Error("Unknown message type received", "type", m.Type)
		return false
	}
}
