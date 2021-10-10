package chain

import "math/big"

type ERC20 interface {
	Name() string
	Symbol() string
	Decimals() uint8
	TotalSupply() big.Int
	BalanceOf() big.Int
	Transfer(string, big.Int) bool
	TransferFrom(string, string, big.Int) bool
	Approve(string, big.Int) bool
	Allowance(string, string) big.Int
	EventTransfer(string, string, big.Int)
	EventApproval(string, string, big.Int)
}
