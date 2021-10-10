package doge

var (
	network = "mainnet"
)

func SetNetwork(name string) {
	network = name
}

func isMainNet() bool {
	return network == "mainnet"
}

func getVer(name string) uint8 {
	switch name {
	case "P2PKH":
		if isMainNet() {
			return 0x1E
		}
		return 0x1E

	case "WIF":
		if isMainNet() {
			return 0x9E
		}
		return 0x9E
	}

	return 0x00
}
