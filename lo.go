package utee

import (
	"cmp"
	"sort"

	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

/*
 wrap and extend github.com/samber/lo
*/

// Clone returns a shallow copy of the collection.
func Clone[T any](collection []T) []T {
	return append(collection[:0:0], collection...)
}

// Map lo.Map without loop index
func Map[T any, R any](collection []T, fn func(T) R) []R {
	return lo.Map(collection, func(t T, _ int) R { return fn(t) })
}

// Filter lo.Filter  without loop index
func Filter[T any](collection []T, fn func(T) bool) []T {
	return lo.Filter(collection, func(t T, _ int) bool { return fn(t) })
}

// Reject lo.Reject  without loop index
func Reject[T any](collection []T, fn func(T) bool) []T {
	return lo.Reject(collection, func(t T, _ int) bool { return fn(t) })
}

// FilterMap lo.FilterMap  without loop index
func FilterMap[T any, R any](collection []T, fn func(T) (R, bool)) []R {
	return lo.FilterMap(collection, func(t T, _ int) (R, bool) { return fn(t) })
}

// FlatMap lo.FlatMap  without loop index
func FlatMap[T any, R any](collection []T, fn func(T) []R) []R {
	return lo.FlatMap(collection, func(t T, _ int) []R { return fn(t) })
}

// Shuffle return a shuffled copy of  collection
func Shuffle[T any](collection []T) []T {
	out := Clone(collection)
	mutable.Shuffle(out)

	return out
}

// Reverse return a reversed copy of collection
func Reverse[T any](collection []T) []T {
	out := Clone(collection)
	mutable.Reverse(out)

	return out
}

// OrderBy order by fn, return ordered copy of slice
func OrderBy[T any, R cmp.Ordered](l []T, fn func(T) R) []T {
	out := Clone(l)
	sort.Slice(out, func(i, j int) bool {
		return fn(out[i]) < fn(out[j])
	})

	return out
}

// OrderByDescending order by fn descending, return ordered copy of slice
func OrderByDescending[T any, R cmp.Ordered](l []T, fn func(T) R) []T {
	out := Clone(l)
	sort.Slice(out, func(i, j int) bool {
		return fn(out[j]) < fn(out[i])
	})

	return out
}
