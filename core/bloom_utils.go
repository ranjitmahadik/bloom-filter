package core

func setBit(arr []byte, val uint64) {
	idx, offset := val/8, val%8
	if idx >= uint64(len(arr)) {
		return
	}
	arr[idx] |= (1 << offset)
}

func isBitSet(arr []byte, val uint64) bool {
	idx, offset := val/8, val%8
	if idx >= uint64(len(arr)) {
		return false
	}
	return (arr[idx]&(1<<offset)&(1<<offset) > 0)
}
