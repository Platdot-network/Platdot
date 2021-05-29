<p align="center"><a href=""><img width="400" title="Platdot" src='https://cdn.jsdelivr.net/gh/rjman-self/resources/assets/platdot-logo.PNG' /></a></p>

# Platdot

[English](../README.md) | 简体中文

![build](https://img.shields.io/badge/build-passing-{})    ![test](https://img.shields.io/badge/test-passing-{})   ![release](https://img.shields.io/badge/release-v1.0.0-E6007A)    ![license](https://img.shields.io/badge/License-Apache%202.0-blue?logo=apache&style=flat-square)

[communicate with us](https://matrix.to/#/#platdot-faucet:matrix.org?via=matrix.org)

## Platdot跨链桥

`Platdot` 是基于 [ChainBridge](https://github.com/ChainSafe/ChainBridge) 开发的一个跨链项目，为Platon提供了Polkadot的跨链桥，实现Pdot的发行、赎回、转账功能。

## 总览

目前，Platdot基于受信任的联盟模型下运行，用户可以非常低的手续费完成抵押发行和赎回操作。现在正处于测试阶段，实现了KSM（Kusama）、AKSM（Alaya）的流通，[UI](https://github.com/Platdot-network/Platdot-UI)界面如下所示：

![Platdot-ui](https://camo.githubusercontent.com/ea47333872a3c8b9ec4eda5aff4412d305f05df7e781c4a6947f3fa617ed9396/68747470733a2f2f6674702e626d702e6f76682f696d67732f323032312f30332f323338386337613738353631383734372e706e67)

在PlatON网络上，EVM的智能合约能够在接收事务时实现自定的处理行为，如发行、销毁新的资产。在Polkadot网络上，multisig-pallet提供了多重签名功能，Platdot据此进行设计完成抵押和赎回操作。例如，Polkadot网络上锁定DOT资产，在EVM上执行合约能铸造并发行PDOT资产，同样地，在EVM上执行合约能销毁PDOT资产，并从Polkadot的多重签名地址赎回DOT资产。

![Platdot-overview](Platdot-overview.png)

## 安装

### 依赖项

- Make sure the Golang environment is installed

### 构建

`git clone https://github.com/RJman-self/Platdot.git`

`make build`: Builds `platdot` in `./build`.

**or**

`make install`: Uses `go install` to add `platdot` to your `GOBIN`.

## 启动Platdot

查阅下方 `GitHub Wiki` 文档.

[作为Relayer启动Platdot](https://github.com/Platdot-network/Platdot/wiki/Start-Platdot-as-a-relayer)

## License

The project is released under the terms of the `Apache v2`.
