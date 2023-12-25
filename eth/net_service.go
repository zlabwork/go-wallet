package eth

// docs: https://goethereumbook.org/zh

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
	"math/big"
	"regexp"
)

type Handle struct {
	cli     *ethclient.Client
	chainID *big.Int
}

// NewServiceHandle
// dsn demo
// http://127.0.0.1:8545
// wss://mainnet.infura.io/ws/v3/xxxxxxxx
// https://mainnet.infura.io/v3/xxxxxxxx
func NewServiceHandle(dsn string) (*Handle, error) {
	c, err := ethclient.Dial(dsn)
	if err != nil {
		return nil, err
	}
	return &Handle{cli: c}, nil
}

// GetClient
// 获取连接
func (h *Handle) GetClient() *ethclient.Client {
	return h.cli
}

// SetChainID
// set chainId
func (h *Handle) SetChainID(id int64) {
	h.chainID = big.NewInt(id)
}

// GetChainID
// get chainId
func (h *Handle) GetChainID() *big.Int {

	if h.chainID != nil {
		return h.chainID
	}
	if h.cli == nil {
		return nil
	}

	id, err := h.cli.NetworkID(context.Background())
	if err != nil {
		return nil
	}
	h.chainID = id
	return h.chainID
}

// GetBalance
// 获取额度
func (h *Handle) GetBalance(address string) (*big.Int, error) {
	if h.cli == nil {
		return nil, errors.New("server is not connected")
	}
	account := common.HexToAddress(address)
	balance, err := h.cli.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// IsValidAddress
// 校验地址
func (h *Handle) IsValidAddress(address string) bool {
	reg := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return reg.MatchString(address)
}

// IsContract
// 检测是否为合约地址
func (h *Handle) IsContract(address string) (bool, error) {
	if h.cli == nil {
		return false, errors.New("server is not connected")
	}
	addr := common.HexToAddress(address)
	byteCode, err := h.cli.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return false, err
	}
	return len(byteCode) > 0, nil
}

// TransferUsePriKey
// 交易 - 使用指定私钥
func (h *Handle) TransferUsePriKey(priKey []byte, toAddr string, wei *big.Int) (txHash string, err error) {
	k, err := crypto.ToECDSA(priKey)
	if err != nil {
		return "", err
	}
	return h.transferViaPriKey(k, toAddr, wei)
}

func (h *Handle) transferViaPriKey(priKey *ecdsa.PrivateKey, toAddr string, wei *big.Int) (txHash string, err error) {
	if h.cli == nil {
		return "", errors.New("server is not connected")
	}

	address := crypto.PubkeyToAddress(priKey.PublicKey)

	from := address
	to := common.HexToAddress(toAddr)

	nonce, err := h.cli.PendingNonceAt(context.Background(), from)
	if err != nil {
		return "", err
	}

	// gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	gasPrice, err := h.cli.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	gasLimit := uint64(21000) // ETH转账的燃气应设上限为21000单位。

	tx := types.NewTransaction(nonce, to, wei, gasLimit, gasPrice, nil)

	chainID := h.GetChainID()

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), priKey)
	if err != nil {
		return "", err
	}

	// SendTransaction
	err = h.cli.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// CreateTxData
// 生成原始交易 - 简单
func (h *Handle) CreateTxData(priKey []byte, toAddr string, wei *big.Int) (rawTX string, err error) {

	if h.cli == nil {
		return "", errors.New("server is not connected")
	}
	client := h.cli

	// 1. gasLimit gasPrice
	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// 2. chainId
	chainID := h.GetChainID()

	// 3. return
	var data []byte

	return h.CreateTxDataAdvanced(priKey, toAddr, wei, gasLimit, gasPrice, chainID, data)
}

// CreateTxDataAdvanced
// 生成原始交易 - 高级
func (h *Handle) CreateTxDataAdvanced(priKey []byte, toAddr string, wei *big.Int, gasLimit uint64, gasPrice *big.Int, chainID *big.Int, data []byte) (rawTX string, err error) {

	if h.cli == nil {
		return "", errors.New("server is not connected")
	}
	client := h.cli

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
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    wei,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// FIXME :: this is a bug, need to be fixed
	log.Println(signedTx)
	return "", fmt.Errorf("error, needs to fix")

	//ts := types.Transactions{signedTx}
	//rawTxBytes := ts.GetRlp(0)
	//rawTxHex := hex.EncodeToString(rawTxBytes)
	//
	//return rawTxHex, nil
}

// SendRawTX
// 发送原始交易
func (h *Handle) SendRawTX(rawTx string) (txHash string, err error) {

	if h.cli == nil {
		return "", errors.New("server is not connected")
	}

	rawTxBytes, err := hex.DecodeString(rawTx)

	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	if err = h.cli.SendTransaction(context.Background(), tx); err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

// CreateSign
// 签名
func (h *Handle) CreateSign(priKey []byte, data []byte) (signature string, err error) {
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

// GetBlockHeader
// 查询区块头 - number = nil 查询最新区块的头信息
func (h *Handle) GetBlockHeader(number *big.Int) (*types.Header, error) {

	header, err := h.cli.HeaderByNumber(context.Background(), number)
	if err != nil {
		return nil, err
	}

	return header, nil
}

// GetBlock
// 查询区块 - number = nil 查询最新区块
func (h *Handle) GetBlock(number *big.Int) (*types.Block, error) {

	block, err := h.cli.BlockByNumber(context.Background(), number)
	if err != nil {
		return nil, err
	}

	return block, nil
}
