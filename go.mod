module github.com/Platdot-network/Platdot

go 1.15

replace github.com/centrifuge/go-substrate-rpc-client/v3 v3.0.2 => github.com/chainx-org/go-substrate-rpc-client/v3 v3.1.1

require (
	github.com/ChainSafe/log15 v1.0.0
	github.com/Platdot-Network/substrate-go v1.6.7
	github.com/centrifuge/go-substrate-rpc-client/v3 v3.0.2
	github.com/hacpy/chainbridge-substrate-events v1.0.0
	github.com/hacpy/go-ethereum v1.14.1
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/prometheus/client_golang v1.4.1
	github.com/rjman-ljm/go-substrate-crypto v1.0.0
	github.com/rjman-ljm/platdot-utils v1.6.5
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)
