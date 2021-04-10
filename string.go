package utee

// SplitStringSlice split string slice into chunks
func SplitStringSlice(src []string, chunkSize int) [][]string {
	var out [][]string
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
