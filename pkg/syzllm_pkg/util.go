package syzllm_pkg

import (
	"fmt"
	"math/rand"
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
