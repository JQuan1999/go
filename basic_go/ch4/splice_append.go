package ch4

func SpliceAppend(splice, data []byte) []byte {
	n := len(splice)
	m := n + len(data)
	if m > cap(splice) {
		new_splice := make([]byte, m)
		copy(new_splice, splice)
		splice = new_splice
	}
	copy(splice[n:m], data)
	return splice
}
