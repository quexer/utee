package utee

// Assure is a helper function to panic if error is not nil. otherwise, it returns the value.
func Assure[T any](v T, err error) T {
	Chk(err) // fail fast on critical fault
	return v
}

// Chk checks the error and panics if it is not nil.
func Chk(err error) {
	if err != nil {
		panic(err)
	}
}
