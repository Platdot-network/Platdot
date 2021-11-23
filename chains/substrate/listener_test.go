package substrate

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	subtest "github.com/Platdot-network/Platdot/shared/substrate/testing"
	"github.com/Platdot-Network/go-substrate-rpc-client/v3/types"
	"github.com/rjman-ljm/platdot-utils/msg"
)

const ListenerTimeout = time.Second * 30

type mockRouter struct {
	msgs chan msg.Message
}

func (r *mockRouter) Send(message msg.Message) error {
	r.msgs <- message
	return nil
}

func verifyResultingMessage(t *testing.T, r *mockRouter, sysErr chan error, expected msg.Message) {
	// Verify message
	select {
	case err := <-sysErr:
		t.Fatalf("System Error: %s", err)
	case m := <-r.msgs:
		if err := compareMessage(expected, m); err != nil {
			t.Fatal(err)
		}
	case <-time.After(ListenerTimeout):
		t.Fatalf("test timed out")

	}
}

func compareMessage(expected, actual msg.Message) error {
	if !reflect.DeepEqual(expected, actual) {
		if !reflect.DeepEqual(expected.Source, actual.Source) {
			return fmt.Errorf("Source doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.Source, actual.Source)
		} else if !reflect.DeepEqual(expected.Destination, actual.Destination) {
			return fmt.Errorf("Destination doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.Destination, actual.Destination)
		} else if !reflect.DeepEqual(expected.DepositNonce, actual.DepositNonce) {
			return fmt.Errorf("Deposit nonce doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.DepositNonce, actual.DepositNonce)
		} else if !reflect.DeepEqual(expected.Payload, actual.Payload) {
			return fmt.Errorf("Payload doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.Payload, actual.Payload)
		}
	}
	return nil
}

func Test_FungibleTransferEvent(t *testing.T) {
	// Construct our expected message
	var rId msg.ResourceId
	subtest.QueryConst(t, context.client, "Example", "NativeTokenId", &rId)
	amount := big.NewInt(1000000)
	recipient := BobKey.PublicKey
	context.latestOutNonce = context.latestOutNonce + 1
	expected := msg.NewFungibleTransfer(ThisChain, ForeignChain, context.latestOutNonce, amount, rId, recipient)

	subtest.InitiateNativeTransfer(t, context.client, types.NewU128(*amount), recipient, ForeignChain)

	verifyResultingMessage(t, context.router, context.lSysErr, expected)
}

