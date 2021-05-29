# Platdot-Overview

## Platon 背景

`Platon`是一条在以太坊的基础上进行二次开发的公链，用了自己研发bpft共识替代了以太坊的pow共识，拥有evm和wasm双虚拟机。Platon的重心在隐私计算，他们认为数据是驱动人类社会未来发展的石油。

## ChainBridge跨链流程

![ChainBridgeDesign](https://cdn.jsdelivr.net/gh/rjman-self/resources/assets/ChainBridgeDesign.png)

ChainBridge适用条件

+ 在EVM上部署具有多签模块的 [chainbridge-solidity](https://github.com/ChainSafe/chainbridge-solidity) 智能合约
+ 能在基于Substrate开发的链上添加具有多签模块的 [chainbridge-substrate](https://github.com/ChainSafe/chainbridge-substrate) pallet

### Relay的功能

![ChainBridge](https://cdn.jsdelivr.net/gh/rjman-self/resources/assets/20210330101928.png)

ChainA的代币流通到ChainB：
`ChainA Tx` -> `Deposit` -> `Event` -> `Message` -> `parseMessage` -> `ChainB Tx`

+ 用户在ChainA上发起`抵押交易Tx`(Deposit)，与contract（pallet）交互产生`DepositEvent`
+ ChainA的`Listener`持续监听，从存储中检索到Deposit`Event`，构造`Message`转发给对应ChainB的`Writer`
+ Writer解析Message参数，与contract（pallet）交互构造`Proposal提案`
+ Relayer对Proposal进行投票，通过后执行多签交易

#### chainbridge 智能合约和相应pallet的功能

1. 管理relayer、设定投票阈值等功能性活动
2. 构造Proposal，通知Relayer进行投票（emit）
3. 执行Proposal,完成多签交易

#### Listener的功能

1. 监听链上的跨链事件(合约发出的DepositEvent)
2. 构造Message转发给Writer

#### Writer的功能

1. 解析来自Listener的消息
2. 与合约交互创建Proposal
3. 为存在的Proposal进行投票

## Platdot跨链流程

> Platdot在ChainBridge的基础上进行修改，适配了目标链无法自行部署pallet的跨链需求

![PlatdotDesign](https://cdn.jsdelivr.net/gh/rjman-self/resources/assets/PlatdotDesign.png)

Platdot适用条件

+ 在EVM上部署具有多签模块的chainbridge-solidity智能合约
+ 基于Substrate开发的、具有`multisign pallet`的链，如Polkadot / Kusama / Westend

![Platdot](https://cdn.jsdelivr.net/gh/rjman-self/resources/assets/Platdot.jpg)

### 主要变动

#### Listener的变动

Listener监听`区块的数据`，从区块的交易中检索出所需的交易数据，构造参数发给Platon的Writer

#### Writer的变动

获取`链上`已存在的`多签交易`，构造相应的参数，新建/批准多签交易

+ RPC调用获取块数据
+ 解析数据取得需要的参数