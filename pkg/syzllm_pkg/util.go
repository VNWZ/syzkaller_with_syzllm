package syzllm_pkg

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var SyzllmProbabilityFuzzer = 0.35
var MutationSelectionRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type SyscallRequestData struct {
	Syscalls []string
}

type SyzLLMResponse struct {
	State   int
	Syscall string
}

func InsertItem[T any](slice []T, item T, pos int) ([]T, error) {
	if pos < 0 || pos > len(slice) {
		return nil, fmt.Errorf("invalid position %d; must be between 0 and %d", pos, len(slice))
	}

	result := make([]T, len(slice)+1)

	copy(result[:pos], slice[:pos])

	result[pos] = item

	copy(result[pos+1:], slice[pos:])

	return result, nil
}

func AssertSlicesAreEqual(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func SortedMapKeysDesc(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		numI, err := strconv.Atoi(keys[i])
		if err != nil {
			panic("invalid numeric key: " + keys[i])
		}
		numJ, err := strconv.Atoi(keys[j])
		if err != nil {
			panic("invalid numeric key: " + keys[j])
		}
		return numI > numJ
	})
	return keys
}

func enlargeSlice(slice []string, pos int) ([]string, error) {
	if pos < 0 || pos > len(slice) {
		return nil, fmt.Errorf("invalid position %d", pos)
	}
	return append(slice[:pos], append([]string{""}, slice[pos:]...)...), nil
}
