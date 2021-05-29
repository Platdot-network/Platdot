<p align="center"><a href=""><img width="400" title="Platdot" src='https://cdn.jsdelivr.net/gh/rjman-self/resources/assets/platdot-logo.PNG' /></a></p>

# Platdot

English | [简体中文](./docs/README_CN.md)

![build](https://img.shields.io/badge/build-passing-{})    ![test](https://img.shields.io/badge/test-passing-{})   ![release](https://img.shields.io/badge/release-v1.0.0-E6007A)    ![license](https://img.shields.io/badge/License-Apache%202.0-blue?logo=apache&style=flat-square)

[communicate with us](https://matrix.to/#/#platdot-faucet:matrix.org?via=matrix.org)

## The cross-chain Bridge

`Platdot` is a cross-chain bridge based on [ChainBridge](https://github.com/ChainSafe/ChainBridge), it provides `Polkadot` cross-chain bridge for `Platon` to realize the functions of PDOT **issuance**, **redemption** and **transfer**.

## Overview

Currently, Platdot operates under a trusted federation model, and users can complete mortgage issuance and redemption operations at a very low handling fee. It is now in the testing phase and has realized the circulation of KSM and AKSM.The [UI](https://github.com/Platdot-network/Platdot-UI) page looks like this:

![Platdot-ui](https://camo.githubusercontent.com/ea47333872a3c8b9ec4eda5aff4412d305f05df7e781c4a6947f3fa617ed9396/68747470733a2f2f6674702e626d702e6f76682f696d67732f323032312f30332f323338386337613738353631383734372e706e67)

## Features

On the PlatON network, EVM's smart contract can implement custom processing behaviors when receiving transactions, such as issuing and destroying new token. On the Polkadot network, multisig-pallet provides a multi-signature function, and Platdot designs accordingly to complete mortgage and redemption. For example, locking DOT assets on the Polkadot network and executing contracts on the EVM can mint and issue PDOT assets. Similarly, executing contracts on the EVM can destroy PDOT assets and redeem DOT assets from Polkadot's multi-signature addresses.

![Platdot-overview](./docs/Platdot-overview.png)

## Installation

### Prerequisites

- Make sure the `Go` environment is installed

### Building

`git clone https://github.com/RJman-self/Platdot.git`

`make build`: Builds `platdot` in `./build`.

**or**

`make install`: Uses `go install` to add `platdot` to your `GOBIN`.

## Getting Start

Documentations are now moved to `GitHub Wiki`.

[Start Platdot as a relayer](https://github.com/Platdot-network/Platdot/wiki/Start-Platdot-as-a-relayer)

## License

The project is released under the terms of the `Apache v2`.
