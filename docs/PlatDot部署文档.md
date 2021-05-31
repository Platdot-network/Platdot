## 1，环境要求

- Go
- GCC
- Platdot
- Ubuntu 18.04/20.04
## 2，配置运行环境

### 2.1，安装 Go:

2.1.1，下载Go语言环境

```shell
wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
```

2.1.2，将下载的压缩包解压到/usr/local，在/usr/local/go中创建Go目录。

> 重要提示：此步骤将在解压之前删除位于/usr/local/go的先前安装(如果有)，请先备份所有数据，然后再继续。

```shell
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz
```

2.1.3，添加 `/usr/local/go/bin` 到 `$PATH`

```shell
export PATH=$PATH:/usr/local/go/bin
```
> 立即应用更改，只需生效配置，如`source~/.bashrc`

2.1.4，运行下列命令验证go环境是否安装完成

```shell
go version
```

### 2.2，安装 GCC

```shell
sudo apt update
sudo apt install build-essential
gcc --version
```

### 2.3，从源安装PlatDot

2.3.1，下载PlatDot源码

```sh
sudo apt install git
git clone https://github.com/RJman-self/Platdot
```

2.3.2，编译PlatDot

```shell
cd Platdot && make build
```
编译完成后能在`./build/`目录下找到PlatDot可执行文件

2.3.3，验证是否编译成功

```shell
./build/platdot --version
```

## 3，启动PlatDot

在启动PlatDot前，你需要连接见证人账户，并设定维护的多签地址

### 3.1，创建多签账户

3.1.1，创建多签账户
![Add Multisig](https://cdn.jsdelivr.net/gh/hacpy/PictureBed@master/Platdot/1617253951015-1617253951012.png)

3.1.2，给见证人账户添加多签地址的权限

![Add your account to multisig account](https://cdn.jsdelivr.net/gh/hacpy/PictureBed@master/Platdot/1617254071745-1617254071742.png)

3.1.3，成功创建多签地址
![Create multisig account successful](https://cdn.jsdelivr.net/gh/hacpy/PictureBed@master/Platdot/1617254470399-1617254470390.png)

### 3.2，参考配置文件

启动PlatDot需要一个配置文件，其中包含有关您的帐户和其他中继器的信息，可以参考：

```json
{
	"chains": [
    {
      "name": "kusama_ksm",
      "type": "substrate",
      "id": "1",
      "endpoint": "connect-your-node",
      "from": "5CHwt8bFyDLC3MyzPQugmmxZTGjShBW2kFMWiC2kSL5TuJxd",
      "opts": {
        "MultiSignAddress": "0x83b0e4664507e7072dd2b30e9c5f68a708979e741a965048fb1ccbc61bd331f5",
        "TotalRelayer": "5",
        "CurrentRelayerNumber": "1",
        "MultiSignThreshold": "3",
        "OtherRelayer1": "0x50a80eb26a7fb43ff4f84ead705fc61c1d4074112e53f781a6b03c0c7504f663",
        "OtherRelayer2": "0x923eeef27b93315c97e63e0c1284b7433ffbc413a58da0626a63955a48586075",
        "OtherRelayer3": "0xa45a0ddd81da79f65cbcfeefc8e62382b1f56ccbbdd9533f77cdc49172cca33d",
        "OtherRelayer4": "0xe6c2b6c4a5d3a770814f3ebe99893d1bb66e8f0d086a2badfcbb481b043ada1a",
        "MaxWeight": "22698000000",
        "DestId": "2",
        "ResourceId": "0x0000000000000000000000000000000000000000000000000000000000000000"
      }
    },
	{
      "name": "alaya",
      "type": "ethereum",
      "id": "2",
      "endpoint": "http://47.110.34.31:6789",
      "from": "atp18hqda4eajphkfarxaa2rutc5dwdwx9z5vy2nmh",
      "latestBlock": "true",
      "opts": {
        "bridge": "atp1emxqzwmz0nv5pxk3h9e2dp3p6djfkqwn4v05zk",
        "erc20Handler": "atp15nqwyjpffntmgg05aq6u7frdvy60qnm82007q5",
        "http": "true",
        "prefix": "atp",
        "networkId": "201030"
      }
    },
    {
      "name": "platon",
      "type": "ethereum",
      "id": "4",
      "endpoint": "http://47.241.98.219:6789",
      "from": "lat18hqda4eajphkfarxaa2rutc5dwdwx9z54jutyc",
      "latestBlock": "true",
      "opts": {
        "bridge": "lat13kpqglnd5xl699smjulk64v048ku7d50p3yntw",
        "erc20Handler": "lat18svtj54uzxpxunu0q63fsenyy66skz2eaw4lz3",
        "http": "true",
        "prefix": "lat",
        "networkId": "210309"
      }
    }
  ]
}
```

### 3.3，设定配置文件

From地址为需要连接的见证人账户地址

```json
{
    "name": "alaya",
    "...": "...",
    "from": "your-alaya-account-address"
},
{
    "name": "kuama",
    "...": "...",
	"from": "your-kusama-account-address"
}
```

> 以下Substrate网络的配置信息全都要使用账户的公钥

如果没有其它地址转换工具，可以使用[Polka Address Conversion](https://polkadot.subscan.io/tools/ss58_transform)

在上述操作后我们成功创建了多签账户，在[polkadot app](https://polkadot.js.org/apps/#/accounts) 点击多签账户名称即可看到详细信息，接着我们按照这个信息完成我们的PlatDot配置文件。

![Open multisig account detail](https://cdn.jsdelivr.net/gh/hacpy/PictureBed@master/Platdot/1617257673522-1617257673520.png)

![Multisig account details](https://cdn.jsdelivr.net/gh/hacpy/PictureBed@master/Platdot/1617257781015-1617257781011.png)



如图所示，最上方即为多签地址，`MultiSignAddress`填入如下地址：

> 此处地址需要填写公钥

```json
{
	"name": "kusama",
    "...": "...",
	"MultiSignAddress": "multisig-address-public-key",
}
```

`TotalRelayer` 和`MultiSignThreshold` 分别是见证人的数量和多签交易完成的阈值，`CurrentRelayerNumber`是设定的见证人次序。如上图所示，如果你是见证人Bob,那么可以设定``TotalRelayer` = 5，`MultiSignThreshold` = 3，`CurrentRelayerNumber` is 2.

```json
{
    "name": "kusama",
    "...": "...",
    "TotalRelayer": "5",
    "CurrentRelayerNumber": "2",
    "MultiSignThreshold": "3",
}
```

`OtherRelayer` 是除去当前见证人账户外的其他见证人，如上图所示，分别是DAVE,CHARLIE,ALICE,EVE,依序将地址填入如下字段：

> 此处地址需要填写公钥

```json
{
    "name": "kusama",
    "...": "...",
    "OtherRelayer1": "DAVE-address-public-key",
    "OtherRelayer2": "CHARLIE-address-public-key",
    "OtherRelayer3": "ALICE-address-public-key",
    "OtherRelayer4": "EVE-address-public-key",
}
```

## 4，运行PlatDot

4.1，导入你的Alaya私钥，生成key文件

```shell
./build/platdot accounts import --privateKey "your-private-key" --secp256k1
```

4.2，导入你的Kusama网络私钥，生成key文件

```shell
./build/platdot accounts import --privateKey "your-private-key" --sr25519
```

4.3，根据配置文件运行PlatDot

```sh
./build/platdot --config config.json
```

