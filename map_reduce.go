package utee

import (
	"github.com/kevwan/mapreduce/v2"
)

// dummyReducer in case map reduce don't need reducer
var dummyReducer = func(pipe <-chan any, cancel func(error)) {}

// Finish concurrent run fn, return error if any of fn return error
// similar to mapreduce.Finish, but with concurrent limit
func Finish(concurrent int, fn ...func() error) error {
	if len(fn) == 0 {
		return nil
	}

	gen := func(source chan<- func() error) {
		for _, v := range fn {
			source <- v
		}
	}

	mapper := func(fn func() error, writer mapreduce.Writer[any], cancel func(error)) {
		if err := fn(); err != nil {
			cancel(err)
		}
	}

	return mapreduce.MapReduceVoid(gen, mapper, dummyReducer, mapreduce.WithWorkers(concurrent))
}
