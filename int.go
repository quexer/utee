package utee

// SplitIntSlice split int slice into chunks
func SplitIntSlice(src []int, chunkSize int) [][]int {
	var out [][]int
	for {
		if len(src) == 0 {
			break
		}
		if len(src) < chunkSize {
			chunkSize = len(src)
		}
		out = append(out, src[0:chunkSize])
		src = src[chunkSize:]
	}
	return out
}

// Min returns the smallest parameter
func Min(n ...int) int {
	var out = n[0]
	for _, v := range n {
		if v < out {
			out = v
		}
	}
	return out
}

// Max returns the biggest parameter
func Max(n ...int) int {
	var out = n[0]
	for _, v := range n {
		if v > out {
			out = v
		}
	}
	return out
}
