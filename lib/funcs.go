package lib

type Option[T any] func(*T)

func ApplyOptions[T any](target *T, opts ...Option[T]) {
	for _, opt := range opts {
		opt(target)
	}
}

func Reverse[T any](slice []T) []T {
	n := len(slice)
	rev := make([]T, n)
	for i := len(slice) - 1; i >= 0; i-- {
		rev = append(rev, slice[i])
	}
	return rev
}

func ReverseInPlace[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
