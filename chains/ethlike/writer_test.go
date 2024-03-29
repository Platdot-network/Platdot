package ethlike

import (
	"context"
	"github.com/ChainSafe/log15"
	"github.com/Platdot-network/Platdot/bindings/Bridge"
	"github.com/Platdot-network/Platdot/chains/chainset"
	utils "github.com/Platdot-network/Platdot/shared/ethlike"
	ethtest "github.com/Platdot-network/Platdot/shared/ethlike/testing"
	eth "github.com/hacpy/go-ethereum"
	"github.com/hacpy/go-ethereum/common"
	ethtypes "github.com/hacpy/go-ethereum/core/types"
	ethcrypto "github.com/hacpy/go-ethereum/crypto"
	"github.com/rjman-ljm/platdot-utils/msg"
	"math/big"
	"testing"
	"time"
)

func createWriters(t *testing.T, client *utils.Client, contracts *utils.DeployedContracts) (*writer, *writer, func(), func(), chan error, chan error) {
	latestBlock := ethtest.GetLatestBlock(t, client)
	errA := make(chan error)
	writerA, stopA := createTestWriter(t, createConfig("bob", latestBlock, contracts), errA)
	errB := make(chan error)
	writerB, stopB := createTestWriter(t, createConfig("charlie", latestBlock, contracts), errB)
	return writerA, writerB, stopA, stopB, errA, errB
}

func createTestWriter(t *testing.T, cfg *Config, errs chan<- error) (*writer, func()) {

	conn := newLocalConnection(t, cfg)
	stop := make(chan int)

	bc := chainset.NewChainCore(cfg.name)

	writer := NewWriter(conn, cfg, newTestLogger(cfg.name), *AliceKp, stop, errs, nil, bc)

	bridge, err := Bridge.NewBridge(cfg.bridgeContract, conn.Client())
	if err != nil {
		t.Fatal(err)
	}

	writer.setContract(bridge)

	err = writer.start()
	if err != nil {
		t.Fatal(err)
	}
	return writer, func() { close(stop) }
}

func routeMessageAndWait(t *testing.T, client *utils.Client, alice, bob *writer, m msg.Message, aliceErr, bobErr chan error) {
	// Watch for executed event
	query := eth.FilterQuery{
		FromBlock: big.NewInt(0),
		Addresses: []common.Address{alice.cfg.bridgeContract},
		Topics: [][]common.Hash{
			{utils.ProposalEvent.GetTopic()},
		},
	}

	ch := make(chan ethtypes.Log)
	sub, err := client.Client.SubscribeFilterLogs(context.Background(), query, ch)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	if err != nil {
		log15.Error("Failed to subscribe to finalization event", "err", err)
	}

	// Alice processes the message, then waits to execute
	if ok := alice.ResolveMessage(m); !ok {
		t.Fatal("Alice failed to resolve the message")
	}

	// Now Bob receives the same message and also waits to execute
	if ok := bob.ResolveMessage(m); !ok {
		t.Fatal("Bob failed to resolve the message")
	}

	for {
		select {
		case evt := <-ch:
			sourceId := evt.Topics[1].Big().Uint64()
			depositNonce := evt.Topics[2].Big().Uint64()
			status := uint8(evt.Topics[3].Big().Uint64())

			if m.Source == msg.ChainId(sourceId) &&
				uint64(m.DepositNonce) == depositNonce &&
				utils.IsExecuted(status) {
				return
			}

		case err = <-sub.Err():
			if err != nil {
				t.Fatal(err)
			}
		case err = <-aliceErr:
			t.Fatalf("Fatal error: %s", err)
		case err = <-bobErr:
			t.Fatalf("Fatal error: %s", err)
		case <-time.After(TestTimeout):
			t.Fatal("test timed out")
		}
	}
}

func TestWriter_start_stop(t *testing.T) {
	conn := newLocalConnection(t, aliceTestConfig)
	defer conn.Close()

	stop := make(chan int)

	bc := chainset.NewChainCore(aliceTestConfig.name)
	writer := NewWriter(conn, aliceTestConfig, TestLogger, *AliceKp, stop, nil, nil, bc)

	err := writer.start()
	if err != nil {
		t.Fatal(err)
	}

	// Initiate shutdown
	close(stop)
}

func TestCreateAndExecuteErc20DepositProposal(t *testing.T) {
	client := ethtest.NewClient(t, TestEndpoint, AliceKp)
	contracts := deployTestContracts(t, client, TestChainId)
	writerA, writerB, stopA, stopB, errA, errB := createWriters(t, client, contracts)

	defer stopA()
	defer stopB()
	defer writerA.conn.Close()
	defer writerB.conn.Close()
	erc20Address := ethtest.DeployMintApproveErc20(t, client, contracts.ERC20HandlerAddress, big.NewInt(100))
	ethtest.FundErc20Handler(t, client, contracts.ERC20HandlerAddress, erc20Address, big.NewInt(100))

	// Create initial transfer message
	resourceId := msg.ResourceIdFromSlice(append(common.LeftPadBytes(erc20Address.Bytes(), 31), 0))
	recipient := ethcrypto.PubkeyToAddress(BobKp.PrivateKey().PublicKey)
	amount := big.NewInt(10)
	m := msg.NewFungibleTransfer(1, 0, 0, amount, resourceId, recipient.Bytes())
	ethtest.RegisterResource(t, client, contracts.BridgeAddress, contracts.ERC20HandlerAddress, resourceId, erc20Address)
	// Helpful for debugging
	go ethtest.WatchEvent(client, contracts.BridgeAddress, utils.ProposalEvent)
	go ethtest.WatchEvent(client, contracts.BridgeAddress, utils.ProposalVote)

	routeMessageAndWait(t, client, writerA, writerB, m, errA, errB)

	ethtest.Erc20AssertBalance(t, client, amount, erc20Address, recipient)
}

func TestDuplicateMessage(t *testing.T) {
	client := ethtest.NewClient(t, TestEndpoint, AliceKp)
	contracts := deployTestContracts(t, client, TestChainId)
	writerA, writerB, stopA, stopB, errA, errB := createWriters(t, client, contracts)

	defer stopA()
	defer stopB()
	defer writerA.conn.Close()
	defer writerB.conn.Close()

	erc20Address := ethtest.DeployMintApproveErc20(t, client, contracts.ERC20HandlerAddress, big.NewInt(100))
	ethtest.FundErc20Handler(t, client, contracts.ERC20HandlerAddress, erc20Address, big.NewInt(100))

	// Create initial transfer message
	resourceId := msg.ResourceIdFromSlice(append(common.LeftPadBytes(erc20Address.Bytes(), 31), 0))
	recipient := ethcrypto.PubkeyToAddress(BobKp.PrivateKey().PublicKey)
	amount := big.NewInt(10)
	m := msg.NewFungibleTransfer(1, 0, 10, amount, resourceId, recipient.Bytes())
	ethtest.RegisterResource(t, client, contracts.BridgeAddress, contracts.ERC20HandlerAddress, resourceId, erc20Address)

	data := ConstructErc20ProposalData(m.Payload[0].([]byte), m.Payload[1].([]byte))
	dataHash := utils.Hash(append(contracts.ERC20HandlerAddress.Bytes(), data...))

	// Helpful for debugging
	go ethtest.WatchEvent(client, contracts.BridgeAddress, utils.ProposalEvent)
	go ethtest.WatchEvent(client, contracts.BridgeAddress, utils.ProposalVote)

	// Process initial message
	routeMessageAndWait(t, client, writerA, writerB, m, errA, errB)

	ethtest.Erc20AssertBalance(t, client, amount, erc20Address, recipient)

	// Capture nonces
	nonceAPre, err := writerA.conn.Client().PendingNonceAt(context.Background(), writerA.conn.Keypair().CommonAddress())
	if err != nil {
		t.Fatal(err)
	}
	nonceBPre, err := writerA.conn.Client().PendingNonceAt(context.Background(), writerB.conn.Keypair().CommonAddress())
	if err != nil {
		t.Fatal(err)
	}

	// Try processing the same message again
	if ok := writerA.ResolveMessage(m); ok {
		t.Fatalf("%s should have not voted", writerA.cfg.name)
	}
	if ok := writerB.ResolveMessage(m); ok {
		t.Fatalf("%s should have not voted", writerB.cfg.name)
	}

	// Ensure the votes are recorded
	if !writerA.hasVoted(m.Source, m.DepositNonce, dataHash) {
		t.Fatal("Relayer vote not found on chain")
	}
	if !writerB.hasVoted(m.Source, m.DepositNonce, dataHash) {
		t.Fatal("Relayer vote not found on chain")
	}

	// Capture new nonces
	nonceAPost, err := writerA.conn.Client().PendingNonceAt(context.Background(), writerA.conn.Keypair().CommonAddress())
	if err != nil {
		t.Fatal(err)
	}
	nonceBPost, err := writerA.conn.Client().PendingNonceAt(context.Background(), writerB.conn.Keypair().CommonAddress())
	if err != nil {
		t.Fatal(err)
	}

	// Enssure neither relayer has submitted a tx
	if nonceAPost != nonceAPre {
		t.Error("Writer A nonce has changed.")
	}
	if nonceBPost != nonceBPre {
		t.Error("Writer B nonce has changed.")
	}

}

