// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
/*
The ethereum package contains the logic for interacting with ethereum chains.

There are 3 major components: the connection, the listener, and the writer.
The currently supported transfer types are Fungible (ERC20), Non-Fungible (ERC721), and generic.

Connection

The connection contains the ethereum RPC client and can be accessed by both the writer and listener.

Listener

The listener polls for each new block and looks for deposit events in the bridge contract. If a deposit occurs, the listener will fetch additional information from the handler before constructing a message and forwarding it to the router.

Writer

The writer recieves the message and creates a proposals on-chain. Once a proposal is made, the writer then watches for a finalization event and will attempt to execute the proposal if a matching event occurs. The writer skips over any proposals it has already seen.
*/
package ethlike

import (
	"github.com/ChainSafe/log15"
	bridge "github.com/Platdot-network/Platdot/bindings/Bridge"
	erc20Handler "github.com/Platdot-network/Platdot/bindings/ERC20Handler"
	"github.com/Platdot-network/Platdot/chains/chainset"
	"github.com/Platdot-network/Platdot/config"
	connection "github.com/Platdot-network/Platdot/connections/ethlike"
	"github.com/hacpy/go-ethereum/accounts/abi/bind"
	"github.com/hacpy/go-ethereum/common"
	"github.com/hacpy/go-ethereum/ethclient"
	"github.com/rjman-ljm/platdot-utils/blockstore"
	"github.com/rjman-ljm/platdot-utils/core"
	"github.com/rjman-ljm/platdot-utils/crypto/secp256k1"
	"github.com/rjman-ljm/platdot-utils/keystore"
	metrics "github.com/rjman-ljm/platdot-utils/metrics/types"
	"github.com/rjman-ljm/platdot-utils/msg"
	"math/big"
	"strconv"
)

var _ core.Chain = &Chain{}

var _ Connection = &connection.Connection{}

type Connection interface {
	GetEndPoint() string
	Connect() error
	Reconnect(endpoint string) error
	Keypair() *secp256k1.Keypair
	Opts() *bind.TransactOpts
	CallOpts() *bind.CallOpts
	LockAndUpdateOpts() error
	UnlockOpts()
	Client() *ethclient.Client
	EnsureHasBytecode(address common.Address) error
	LatestBlock() (*big.Int, error)
	WaitForBlock(block *big.Int, delay *big.Int) error
	Close()
}

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     Connection        // THe chains connection
	listener *listener         // The listener of this chain
	writer   *writer           // The writer of the chain
	stop     chan<- int
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than cfg.startBlock, then cfg.startBlock is replaced with the latest known block.
func setupBlockstore(cfg *Config, kp *secp256k1.Keypair) (*blockstore.Blockstore, error) {
	bs, err := blockstore.NewBlockstore(cfg.blockstorePath, cfg.id, kp.Address())
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func InitializeChain(chainCfg *core.ChainConfig, logger log15.Logger, sysErr chan<- error, m *metrics.ChainMetrics) (*Chain, error) {
	// parse config
	cfg, err := parseChainConfig(chainCfg)
	if err != nil {
		return nil, err
	}

	bc := chainset.NewChainCore(cfg.name)

	// set chainId
	networkId, _ := strconv.ParseUint(cfg.networkId, 0, 64)

	// load key
	ethBytes, _ := common.PlatonToEth(cfg.from)
	ethAddress := common.BytesToAddress(ethBytes)
	pwdCache := cfg.keystorePath + "/.cache"
	kpI, err := keystore.KeypairFromAddress(
		ethAddress.String(),
		keystore.EthChain,
		cfg.keystorePath,
		chainCfg.Insecure,
		pwdCache,
		ethAddress.String()[:32],
	)
	if err != nil {
		return nil, err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	// init block store
	bs, err := setupBlockstore(cfg, kp)
	if err != nil {
		return nil, err
	}

	stop := make(chan int)
	conn := connection.NewConnection(networkId, cfg.endpoint[config.InitialEndPointId], cfg.http, kp, logger, cfg.gasLimit, cfg.maxGasPrice, cfg.gasMultiplier)
	err = conn.Connect()
	if err != nil {
		return nil, err
	}
	err = conn.EnsureHasBytecode(cfg.bridgeContract)
	if err != nil {
		return nil, err
	}

	bridgeContract, err := bridge.NewBridge(cfg.bridgeContract, conn.Client())
	if err != nil {
		return nil, err
	}

	erc20HandlerContract, err := erc20Handler.NewERC20Handler(cfg.erc20HandlerContract, conn.Client())
	if err != nil {
		return nil, err
	}

	//if chainCfg.LatestBlock {
	if cfg.startBlock.Uint64() == 0 {
		curr, err := conn.LatestBlock()
		if err != nil {
			return nil, err
		}
		cfg.startBlock = curr
		log15.Info("Start block is newest", "StartBlock", cfg.startBlock, "Chain", cfg.name)
	} else {
		log15.Info("Start block is specified", "StartBlock", cfg.startBlock, "Chain", cfg.name)
	}

	listener := NewListener(conn, cfg, logger, bs, stop, sysErr, m)
	listener.setContracts(bridgeContract, erc20HandlerContract)

	writer := NewWriter(conn, cfg, logger, *kp, stop, sysErr, m, bc)
	writer.setContract(bridgeContract)

	return &Chain{
		cfg:      chainCfg,
		conn:     conn,
		writer:   writer,
		listener: listener,
		stop:     stop,
	}, nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.cfg.Id, c.writer)
	c.listener.setRouter(r)
}

func (c *Chain) Start() error {
	err := c.listener.start()
	if err != nil {
		return err
	}

	err = c.writer.start()
	if err != nil {
		return err
	}

	c.writer.log.Debug("Successfully started chain")
	return nil
}

func (c *Chain) Id() msg.ChainId {
	return c.cfg.Id
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

func (c *Chain) LatestBlock() metrics.LatestBlock {
	return c.listener.latestBlock
}

// Stop signals to any running routines to exit
func (c *Chain) Stop() {
	close(c.stop)
	if c.conn != nil {
		c.conn.Close()
	}
}
