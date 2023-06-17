package slices

import (
	"fmt"
	"strings"
)

func Map[T, V any](list []T, f func(in T) V) (result []V) {
	for _, t := range list {
		result = append(result, f(t))
	}

	return result
}

func Distinct[T comparable](a []T) []T {
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	set := make([]T, 0, len(hash))
	for k := range hash {
		set = append(set, k)
	}

	return set
}

func Revert[T any](slice []T) {
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-i-1] = slice[len(slice)-i-1], slice[i]
	}
}

func Fold[T, V any](list []T, f func(item T, value V) V) (result V) {
	for _, t := range list {
		result = f(t, result)
	}

	return result
}

type Set[T comparable] map[T]bool

func ToSet[T comparable](list []T) map[T]bool {
	m := map[T]bool{}

	for _, t := range list {
		m[t] = true
	}

	return m
}

func Index[T any, K comparable](list []T, key func(T) K) map[K]T {
	m := map[K]T{}

	for _, t := range list {
		m[key(t)] = t
	}

	return m
}

func SQLFilter(column, s string) string {
	return strings.Join(
		Map(strings.Split(s, ";"),
			func(mainPart string) string {
				parts := strings.Split(mainPart, ",")

				parts = Map(parts, func(smallPart string) string {
					sql := column

					if strings.HasPrefix(smallPart, "!") {
						smallPart = strings.TrimPrefix(smallPart, "!")
						sql += " not"
					}

					return sql + " like '%" + smallPart + "%'"
				})

				sql := strings.Join(parts, " or ")

				if len(parts) <= 1 {
					return sql
				}

				return fmt.Sprintf("(%s)", sql)
			},
		), " and ")
}

func MultiTrimPrefix(s string, pp []string) (r string) {
	r = s

	for _, p := range pp {
		r = strings.ReplaceAll(r, p, "")
	}

	return r
}
