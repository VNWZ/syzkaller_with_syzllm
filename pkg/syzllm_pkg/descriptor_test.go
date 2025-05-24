package syzllm_pkg

import "testing"

func TestExtractCallNameWithoutDescriptor(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		wantPanic bool
	}{
		{
			name:     "Valid call with descriptor",
			input:    "write$SyzLLM()",
			expected: "write",
		},
		{
			name:     "Valid call without descriptor",
			input:    "dup$SyzLLM()",
			expected: "dup",
		},
		{
			name:     "Complex call name",
			input:    "epoll_create1$SyzLLM()",
			expected: "epoll_create1",
		},
		{
			name:     "no SyzLLM call with res",
			input:    "r19 = epoll_create1$unix()",
			expected: "epoll_create1",
		},
		{
			name:     "call with res",
			input:    "r19 = epoll_create1$SyzLLM()",
			expected: "epoll_create1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("ExtractCallNameWithoutDescriptor() should have panicked")
					}
				}()
				ExtractCallNameWithoutDescriptor(tt.input)
			} else {
				result := ExtractCallNameWithoutDescriptor(tt.input)
				if result != tt.expected {
					t.Errorf("ExtractCallNameWithoutDescriptor() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestProcessDescriptor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Write with nested sendto",
			input:    "write$SyzLLM(@RSTART@sendto$SyzLLM(@RSTART@socket$SyzLLM(0x1, 0x80001, 0x0)@REND@, &(0x7f0000014000)='\x00'/64, 0x18, 0x4000, 0x0, 0x0)@REND@, &(0x7f000000a000), 0x1)",
			expected: "write$damon_target_ids(@RSTART@sendto$llc(@RSTART@socket$unix(0x1, 0x80001, 0x0)@REND@, &(0x7f0000014000)='\x00'/64, 0x18, 0x4000, 0x0, 0x0)@REND@, &(0x7f000000a000), 0x1)",
		},
		{
			name:     "Write with pipe",
			input:    "write$SyzLLM(@RSTART@pipe$SyzLLM()@REND@, &(0x7f000000a000), 0x1)",
			expected: "write$damon_target_ids(@RSTART@pipe()@REND@, &(0x7f000000a000), 0x1)",
		},
		{
			name:     "Dup",
			input:    "dup$SyzLLM(0x3d)",
			expected: "dup(0x3d)",
		},
		{
			name:     "Epoll_wait with nested epoll_create1",
			input:    "epoll_wait$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, &(0x7f0000059000)=[{}], 0x400, 0x0)",
			expected: "epoll_wait(@RSTART@epoll_create1(0x80000)@REND@, &(0x7f0000059000)=[{}], 0x400, 0x0)",
		},
		{
			name:     "pipe2 without SyzLLM",
			input:    "pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
			expected: "pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
		},
		{
			name:     "pipe without SyzLLM",
			input:    "pipe(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
			expected: "pipe(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessDescriptor(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessDescriptor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessNestedDescriptor(t *testing.T) {
	replacements := map[string]string{
		"socket":     "unix",
		"sendto":     "llc",
		"write":      "damon_target_ids",
		"epoll_ctl":  "EPOLL_CTL_ADD",
		"socketpair": "unix",
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Nested socket call",
			input:    "@RSTART@socket$SyzLLM(0x1, 0x80001, 0x0)@REND@",
			expected: "@RSTART@socket$unix(0x1, 0x80001, 0x0)@REND@",
		},
		{
			name:     "Nested pipe call (no descriptor)",
			input:    "@RSTART@pipe$SyzLLM()@REND@",
			expected: "@RSTART@pipe()@REND@",
		},
		{
			name:     "Multiple nested calls",
			input:    "@RSTART@sendto$SyzLLM(@RSTART@socket$SyzLLM(0x1, 0x80001, 0x0)@REND@, &(0x7f0000014000))@REND@",
			expected: "@RSTART@sendto$llc(@RSTART@socket$unix(0x1, 0x80001, 0x0)@REND@, &(0x7f0000014000))@REND@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessNestedDescriptor(tt.input, replacements)
			if result != tt.expected {
				t.Errorf("ProcessNestedDescriptor() = %v, want %v", result, tt.expected)
			}
		})
	}
}
