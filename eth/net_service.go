package eth

// docs: https://goethereumbook.org/zh

import (
    "context"
    "crypto/ecdsa"
    "encoding/hex"
    "errors"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/rlp"
    "math/big"
    "regexp"
)

type EthHandle struct {
    cli     *ethclient.Client
    chainID *big.Int
}

// dsn demo
// http://127.0.0.1:8545
// wss://mainnet.infura.io/ws/v3/xxxxxxxx
// https://mainnet.infura.io/v3/xxxxxxxx
func NewConnectService(dsn string) (*EthHandle, error) {
    c, err := ethclient.Dial(dsn)
    if err != nil {
        return nil, err
    }
    return &EthHandle{cli: c}, nil
}

// 获取连接
func (e *EthHandle) GetClient() *ethclient.Client {
    return e.cli
}

// set chainId
func (e *EthHandle) SetChainID(id int64) {
    e.chainID = big.NewInt(id)
}

// get chainId
func (e *EthHandle) GetChainID() *big.Int {

    if e.chainID != nil {
        return e.chainID
    }
    if e.cli == nil {
        return nil
    }

    id, err := e.cli.NetworkID(context.Background())
    if err != nil {
        return nil
    }
    e.chainID = id
    return e.chainID
}

// 获取额度
func (e *EthHandle) GetBalance(address string) (*big.Int, error) {
    if e.cli == nil {
        return nil, errors.New("server is not connected")
    }
    account := common.HexToAddress(address)
    balance, err := e.cli.BalanceAt(context.Background(), account, nil)
    if err != nil {
        return nil, err
    }
    return balance, nil
}

// 校验地址
func (e *EthHandle) IsValidAddress(address string) bool {
    re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
    return re.MatchString(address)
}

// 检测是否为合约地址
func (e *EthHandle) IsContract(address string) (bool, error) {
    if e.cli == nil {
        return false, errors.New("server is not connected")
    }
    addr := common.HexToAddress(address)
    byteCode, err := e.cli.CodeAt(context.Background(), addr, nil)
    if err != nil {
        return false, err
    }
    return len(byteCode) > 0, nil
}

// 交易 - 使用指定私钥
func (e *EthHandle) TransferUsePriKey(priKey []byte, toAddr string, wei *big.Int) (txHash string, err error) {
    k, err := crypto.ToECDSA(priKey)
    if err != nil {
        return "", err
    }
    return e.transferViaPriKey(k, toAddr, wei)
}

func (e *EthHandle) transferViaPriKey(priKey *ecdsa.PrivateKey, toAddr string, wei *big.Int) (txHash string, err error) {
    if e.cli == nil {
        return "", errors.New("server is not connected")
    }

    address := crypto.PubkeyToAddress(priKey.PublicKey)

    from := address
    to := common.HexToAddress(toAddr)

    nonce, err := e.cli.PendingNonceAt(context.Background(), from)
    if err != nil {
        return "", err
    }

    // gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
    gasPrice, err := e.cli.SuggestGasPrice(context.Background())
    if err != nil {
        return "", err
    }
    gasLimit := uint64(21000) // ETH转账的燃气应设上限为21000单位。

    tx := types.NewTransaction(nonce, to, wei, gasLimit, gasPrice, nil)

    chainID := e.GetChainID()

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), priKey)
    if err != nil {
        return "", err
    }

    // SendTransaction
    err = e.cli.SendTransaction(context.Background(), signedTx)
    if err != nil {
        return "", err
    }

    return signedTx.Hash().Hex(), nil
}

// 生成原始交易 - 简单
func (e *EthHandle) CreateTxData(priKey []byte, toAddr string, wei *big.Int) (rawTX string, err error) {

    if e.cli == nil {
        return "", errors.New("server is not connected")
    }
    client := e.cli

    // 1. gasLimit gasPrice
    gasLimit := uint64(21000) // in units
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        return "", err
    }

    // 2. chainId
    chainID := e.GetChainID()

    // 3. return
    var data []byte

    return e.CreateTxDataAdvanced(priKey, toAddr, wei, gasLimit, gasPrice, chainID, data)
}

// 生成原始交易 - 高级
func (e *EthHandle) CreateTxDataAdvanced(priKey []byte, toAddr string, wei *big.Int, gasLimit uint64, gasPrice *big.Int, chainID *big.Int, data []byte) (rawTX string, err error) {

    if e.cli == nil {
        return "", errors.New("server is not connected")
    }
    client := e.cli

    // 1. 私钥
    privateKey, err := crypto.ToECDSA(priKey)
    if err != nil {
        return "", err
    }

    // 2. 公钥
    address := crypto.PubkeyToAddress(privateKey.PublicKey)

    // 3. nonce
    fromAddress := address
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        return "", err
    }

    // 4.
    toAddress := common.HexToAddress(toAddr)
    tx := types.NewTransaction(nonce, toAddress, wei, gasLimit, gasPrice, data)

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        return "", err
    }

    ts := types.Transactions{signedTx}
    rawTxBytes := ts.GetRlp(0)
    rawTxHex := hex.EncodeToString(rawTxBytes)

    return rawTxHex, nil
}

// 发送原始交易
func (e *EthHandle) SendRawTX(rawTx string) (txHash string, err error) {

    if e.cli == nil {
        return "", errors.New("server is not connected")
    }

    rawTxBytes, err := hex.DecodeString(rawTx)

    tx := new(types.Transaction)
    rlp.DecodeBytes(rawTxBytes, &tx)

    if err = e.cli.SendTransaction(context.Background(), tx); err != nil {
        return "", err
    }

    return tx.Hash().Hex(), nil
}

// 签名
func (e *EthHandle) CreateSign(priKey []byte, data []byte) (signature string, err error) {
    privateKey, err := crypto.ToECDSA(priKey)
    if err != nil {
        return "", err
    }

    hash := crypto.Keccak256Hash(data)

    sign, err := crypto.Sign(hash.Bytes(), privateKey)
    if err != nil {
        return "", err
    }

    return hexutil.Encode(sign), nil
}

// 查询区块头 - number = nil 查询最新区块的头信息
func (e *EthHandle) GetBlockHeader(number *big.Int) (*types.Header, error) {

    header, err := e.cli.HeaderByNumber(context.Background(), number)
    if err != nil {
        return nil, err
    }

    return header, nil
}

// 查询区块 - number = nil 查询最新区块
func (e *EthHandle) GetBlock(number *big.Int) (*types.Block, error) {

    block, err := e.cli.BlockByNumber(context.Background(), number)
    if err != nil {
        return nil, err
    }

    return block, nil
}
