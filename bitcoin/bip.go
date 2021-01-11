package bitcoin

import (
	"strconv"
	"strings"
)

// @docs https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
// m / purpose' / coin_type' / account' / change / address_index
// format := "m/44'/0'/0'/0/0"
func ParseBipPath(format string) ([]uint32, error) {
	firstHardened := uint32(0x80000000)
	result := make([]uint32, 0)
	arr := strings.Split(format, "/")
	for index, item := range arr {
		if index == 0 {
			continue
		}
		n := uint32(0)
		if strings.Contains(item, "'") {
			i, err := strconv.ParseUint(item[:len(item)-1], 10, 32)
			if err != nil {
				return nil, err
			}
			n = firstHardened + uint32(i)
		} else {
			i, err := strconv.ParseUint(item, 10, 32)
			if err != nil {
				return nil, err
			}
			n = uint32(i)
		}
		result = append(result, n)
	}
	return result, nil
}
