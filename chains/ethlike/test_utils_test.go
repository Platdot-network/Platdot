package ethlike

import (
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/Platdot-network/Platdot/bindings/Bridge"
	"github.com/Platdot-network/Platdot/config"
	connection "github.com/Platdot-network/Platdot/connections/ethlike"
	utils "github.com/Platdot-network/Platdot/shared/ethlike"
	"github.com/hacpy/go-ethereum/common"
	"github.com/rjman-ljm/platdot-utils/keystore"
	"github.com/rjman-ljm/platdot-utils/msg"
	"math/big"
	"testing"
	"time"
)

var TestEndpoint = []string{"ws://localhost:8545"}

var TestLogger = newTestLogger("test")
var TestTimeout = time.Second * 30

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]

var TestRelayerThreshold = big.NewInt(1)
var TestChainId = msg.ChainId(0)

var aliceTestConfig = createConfig("alice", nil, nil)

func createConfig(name string, startBlock *big.Int, contracts *utils.DeployedContracts) *Config {
	cfg := &Config{
		name:                   name,
		id:                     0,
		endpoint:               TestEndpoint,
		from:                   name,
		keystorePath:           "",
		blockstorePath:         "",
		freshStart:             true,
		bridgeContract:         common.Address{},
		erc20HandlerContract:   common.Address{},
		erc721HandlerContract:  common.Address{},
		genericHandlerContract: common.Address{},
		gasLimit:               big.NewInt(DefaultGasLimit),
		maxGasPrice:            big.NewInt(DefaultGasPrice),
		gasMultiplier:          big.NewFloat(DefaultGasMultiplier),
		http:                   false,
		startBlock:             startBlock,
		blockConfirmations:     big.NewInt(3),
	}

	if contracts != nil {
		cfg.bridgeContract = contracts.BridgeAddress
		cfg.erc20HandlerContract = contracts.ERC20HandlerAddress
	}

	return cfg
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(log15.LvlError, tLog.GetHandler()))
	return tLog
}

func newLocalConnection(t *testing.T, cfg *Config) *connection.Connection {
	kp := keystore.TestKeyRing.EthereumKeys[cfg.from]
	conn := connection.NewConnection(
		DefaultNetworkId,
		TestEndpoint[config.InitialEndPointId],
		false,
		kp,
		TestLogger,
		big.NewInt(DefaultGasLimit),
		big.NewInt(DefaultGasPrice),
		big.NewFloat(DefaultGasMultiplier))
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func deployTestContracts(t *testing.T, client *utils.Client, id msg.ChainId) *utils.DeployedContracts {
	contracts, err := utils.DeployContracts(
		client,
		uint8(id),
		TestRelayerThreshold,
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("=======================================================")
	fmt.Printf("Bridge: %s\n", contracts.BridgeAddress.Hex())
	fmt.Printf("Erc20Handler: %s\n", contracts.ERC20HandlerAddress.Hex())
	fmt.Println("========================================================")

	return contracts
}

func createErc20Deposit(
	t *testing.T,
	contract *Bridge.Bridge,
	client *utils.Client,
	rId msg.ResourceId,
	destRecipient common.Address,
	destId msg.ChainId,
	amount *big.Int,
) {

	data := utils.ConstructErc20DepositData(destRecipient.Bytes(), amount)

	// Incrememnt Nonce by one
	client.Opts.Nonce = client.Opts.Nonce.Add(client.Opts.Nonce, big.NewInt(1))
	if _, err := contract.Deposit(
		client.Opts,
		uint8(destId),
		rId,
		data,
	); err != nil {
		t.Fatal(err)
	}
}