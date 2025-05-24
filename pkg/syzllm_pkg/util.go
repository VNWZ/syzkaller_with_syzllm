package syzllm_pkg

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var SyzllmProbabilityFuzzer = 1.0
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
		return keys[i] > keys[j]
	})
	return keys
}
