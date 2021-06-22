package btc

var (
    network = "mainnet" // mainnet, testnet
)

func SetNetwork(name string) {
    network = name
}

func isMainNet() bool {
    return network == "mainnet"
}

// name = P2PKH, P2SH, WIF
func getVer(name string) uint8 {
    switch name {
    case "P2PKH":
        if isMainNet() {
            return 0x00
        }
        return 0x6F

    case "P2SH":
        if isMainNet() {
            return 0x05
        }
        return 0xC4

    case "WIF":
        if isMainNet() {
            return 0x80
        }
        return 0xEF
    }

    return 0x00
}
