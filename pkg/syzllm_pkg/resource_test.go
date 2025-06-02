package syzllm_pkg

import (
	"testing"
)

// TestHasNormalResTag tests the hasNormalResTag function.
func TestHasNormalResTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid normal resource tag",
			input:    "@RSTART@resource@REND@",
			expected: true,
		},
		{
			name:     "Missing RPrefix",
			input:    "resource@REND@",
			expected: false,
		},
		{
			name:     "Missing RSuffix",
			input:    "@RSTART@resource",
			expected: false,
		},
		{
			name:     "No tags",
			input:    "resource",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Partial prefix",
			input:    "@RSTARTresource@REND@",
			expected: false,
		},
		{
			name:     "Multiple tags",
			input:    "@RSTART@resource@REND@@RSTART@another@REND@",
			expected: true,
		},
		{
			name:     "Case sensitivity",
			input:    "@rstart@resource@rend@",
			expected: false,
		},
		{
			name:     "R",
			input:    "pwrite64$SyzLLM(@RSTART@openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)='./file0\\x00', 0x42, 0x180)@REND@, 0x7f7026e75ead, 0x1, 0x1)",
			expected: true,
		},
		{
			name:     "PIPE",
			input:    "read$SyzLLM(@PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f0000055000)=\"\"/0x400, 0x80)",
			expected: false,
		},
		{
			name:     "R+PIPE",
			input:    "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x1, @PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f000003a000)={0x1, 0x10})",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNormalResTag(tt.input)
			if result != tt.expected {
				t.Errorf("hasNormalResTag(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHasPipeTag tests the hasPipeTag function.
func TestHasPipeTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid pipe tag",
			input:    "@PIPESTART@data@PIPEEND@",
			expected: true,
		},
		{
			name:     "Missing PIPEPrefix",
			input:    "data@PIPEEND@",
			expected: false,
		},
		{
			name:     "Missing PIPESuffix",
			input:    "@PIPESTART@data",
			expected: false,
		},
		{
			name:     "No tags",
			input:    "data",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Partial prefix",
			input:    "@PIPESTARTdata@PIPEEND@",
			expected: false,
		},
		{
			name:     "Multiple tags",
			input:    "@PIPESTART@data@PIPEEND@@PIPESTART@more@PIPEEND@",
			expected: true,
		},
		{
			name:     "Case sensitivity",
			input:    "@pipestart@data@pipeend@",
			expected: false,
		},
		{
			name:     "R",
			input:    "pwrite64$SyzLLM(@RSTART@openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)='./file0\\x00', 0x42, 0x180)@REND@, 0x7f7026e75ead, 0x1, 0x1)",
			expected: false,
		},
		{
			name:     "PIPE",
			input:    "read$SyzLLM(@PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f0000055000)=\"\"/0x400, 0x80)",
			expected: true,
		},
		{
			name:     "R+PIPE",
			input:    "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x1, @PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f000003a000)={0x1, 0x10})",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPipeTag(tt.input)
			if result != tt.expected {
				t.Errorf("hasPipeTag(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHasResTag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "R",
			input:    "pwrite64$SyzLLM(@RSTART@openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)='./file0\\x00', 0x42, 0x180)@REND@, 0x7f7026e75ead, 0x1, 0x1)",
			expected: true,
		},
		{
			name:     "PIPE",
			input:    "read$SyzLLM(@PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f0000055000)=\"\"/0x400, 0x80)",
			expected: true,
		},
		{
			name:     "R+PIPE",
			input:    "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x1, @PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f000003a000)={0x1, 0x10})",
			expected: true,
		},
		{
			name:     "Valid pipe tag",
			input:    "@PIPESTART@data@PIPEEND@",
			expected: true,
		},
		{
			name:     "Missing PIPEPrefix",
			input:    "data@PIPEEND@",
			expected: false,
		},
		{
			name:     "Missing PIPESuffix",
			input:    "@PIPESTART@data",
			expected: false,
		},
		{
			name:     "No tags",
			input:    "writev$SyzLLM(0xb, 0x7ffd4b4cbca0, 0x2)",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasResTag(tt.input)
			if result != tt.expected {
				t.Errorf("hasResTag(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountNormalResProvider(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		shouldPanic bool
		expectedMsg string
	}{
		{
			name:        "No Prov - writev",
			input:       "writev$SyzLLM(0xb, 0x7ffd4b4cbca0, 0x2)",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "No Prov - poll",
			input:       "poll$SyzLLM(&(0x7f0000080000)=[{test_call$(&(0x7f9999)=0, 3)}], 0x1, 0x7d0)",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "1 Normal Prov - bpf",
			input:       "r1 = bpf$PROG_LOAD(0x5, &(0x7f00000000c0)={0x11, 0xd, &(0x7f0000000d40)=ANY=[@ANYBLOB=\"180000000000e3ff000000000000000018110000\", @ANYRES32=r0, @ANYBLOB=\"0000000000000000b7080000000000007b8af8ff00000000bfa200000000000007020000f8ffffffb703000008000000b7040000000000008500000001000000850000000700000095\"], &(0x7f0000000240)='GPL\\x00', 0x0, 0x0, 0x0, 0x0, 0x0, '\\x00', 0x0, @fallback, 0xffffffffffffffff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, @void, @value}, 0x90)\n",
			expected:    1,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "1 multi Prov - pipe",
			input:       "pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "panic - has res tag",
			input:       "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x1, @PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f000003a000)={0x1, 0xC})",
			expected:    0,
			shouldPanic: true,
			expectedMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use defer to catch panic
			var panicMsg interface{}
			defer func() {
				if r := recover(); r != nil {
					panicMsg = r
				}
			}()

			// Call the function
			result := countNormalResProvider(tt.input)
			if result != tt.expected {
				t.Errorf("countNormalResProvider(%q) = %v; want %v", tt.input, result, tt.expected)
			}

			// Check panic behavior
			if tt.shouldPanic {
				if panicMsg == nil {
					t.Errorf("assertNoResToParse(%q) did not panic; expected panic", tt.input)
				} else if panicMsg != tt.expectedMsg {
					t.Errorf("assertNoResToParse(%q) panicked with %v; expected %v", tt.input, panicMsg, tt.expectedMsg)
				}
			} else {
				if panicMsg != nil {
					t.Errorf("assertNoResToParse(%q) panicked with %v; expected no panic", tt.input, panicMsg)
				}
			}
		})
	}
}

func TestCountMultiResProvider(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		shouldPanic bool
		expectedMsg string
	}{
		{
			name:        "No Prov - writev",
			input:       "writev$SyzLLM(0xb, 0x7ffd4b4cbca0, 0x2)",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "No Prov - poll",
			input:       "poll$SyzLLM(&(0x7f0000080000)=[{test_call$(&(0x7f9999)=0, 3)}], 0x1, 0x7d0)",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "1 Normal Prov - bpf",
			input:       "r1 = bpf$PROG_LOAD(0x5, &(0x7f00000000c0)={0x11, 0xd, &(0x7f0000000d40)=ANY=[@ANYBLOB=\"180000000000e3ff000000000000000018110000\", @ANYRES32=r0, @ANYBLOB=\"0000000000000000b7080000000000007b8af8ff00000000bfa200000000000007020000f8ffffffb703000008000000b7040000000000008500000001000000850000000700000095\"], &(0x7f0000000240)='GPL\\x00', 0x0, 0x0, 0x0, 0x0, 0x0, '\\x00', 0x0, @fallback, 0xffffffffffffffff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, @void, @value}, 0x90)\n",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "1 multi Prov - pipe",
			input:       "pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})",
			expected:    2,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "panic - has res tag",
			input:       "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x1, @PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f000003a000)={0x1, 0xC})",
			expected:    0,
			shouldPanic: true,
			expectedMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use defer to catch panic
			var panicMsg interface{}
			defer func() {
				if r := recover(); r != nil {
					panicMsg = r
				}
			}()

			// Call the function
			result := countMultiResProvider(tt.input)
			if result != tt.expected {
				t.Errorf("countMultiResProvider(%q) = %v; want %v", tt.input, result, tt.expected)
			}

			// Check panic behavior
			if tt.shouldPanic {
				if panicMsg == nil {
					t.Errorf("assertNoResToParse(%q) did not panic; expected panic", tt.input)
				} else if panicMsg != tt.expectedMsg {
					t.Errorf("assertNoResToParse(%q) panicked with %v; expected %v", tt.input, panicMsg, tt.expectedMsg)
				}
			} else {
				if panicMsg != nil {
					t.Errorf("assertNoResToParse(%q) panicked with %v; expected no panic", tt.input, panicMsg)
				}
			}
		})
	}
}

func TestCountResProvider(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		shouldPanic bool
		expectedMsg string
	}{
		{
			name:        "No Prov - writev",
			input:       "writev$SyzLLM(0xb, 0x7ffd4b4cbca0, 0x2)",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "No Prov - poll",
			input:       "poll$SyzLLM(&(0x7f0000080000)=[{test_call$(&(0x7f9999)=0, 3)}], 0x1, 0x7d0)",
			expected:    0,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "1 Normal Prov - bpf",
			input:       "r1 = bpf$PROG_LOAD(0x5, &(0x7f00000000c0)={0x11, 0xd, &(0x7f0000000d40)=ANY=[@ANYBLOB=\"180000000000e3ff000000000000000018110000\", @ANYRES32=r0, @ANYBLOB=\"0000000000000000b7080000000000007b8af8ff00000000bfa200000000000007020000f8ffffffb703000008000000b7040000000000008500000001000000850000000700000095\"], &(0x7f0000000240)='GPL\\x00', 0x0, 0x0, 0x0, 0x0, 0x0, '\\x00', 0x0, @fallback, 0xffffffffffffffff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, @void, @value}, 0x90)\n",
			expected:    1,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "1 multi Prov - pipe",
			input:       "pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})",
			expected:    2,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "1 normal + 1 multi Prov",
			input:       "r19 = syscall(&(0x7f0000064000)={<r20=>0xffffffffffffffff, <r21=>0xffffffffffffffff})",
			expected:    3,
			shouldPanic: false,
			expectedMsg: "",
		},
		{
			name:        "panic - has res tag",
			input:       "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x1, @PIPESTART@pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})@PIPEEND@, &(0x7f000003a000)={0x1, 0xC})",
			expected:    0,
			shouldPanic: true,
			expectedMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use defer to catch panic
			var panicMsg interface{}
			defer func() {
				if r := recover(); r != nil {
					panicMsg = r
				}
			}()

			// Call the function
			result := countResProvider(tt.input)
			if result != tt.expected {
				t.Errorf("countResProvider(%q) = %v; want %v", tt.input, result, tt.expected)
			}

			// Check panic behavior
			if tt.shouldPanic {
				if panicMsg == nil {
					t.Errorf("assertNoResToParse(%q) did not panic; expected panic", tt.input)
				} else if panicMsg != tt.expectedMsg {
					t.Errorf("assertNoResToParse(%q) panicked with %v; expected %v", tt.input, panicMsg, tt.expectedMsg)
				}
			} else {
				if panicMsg != nil {
					t.Errorf("assertNoResToParse(%q) panicked with %v; expected no panic", tt.input, panicMsg)
				}
			}
		})
	}
}

func TestUpdateResNum(t *testing.T) {
	tests := []struct {
		name     string
		call     string
		count    int
		expected string
	}{
		{
			name:     "no res Prov - writev",
			call:     "writev$SyzLLM(0xb, 0x7ffd4b4cbca0, 0x2)",
			count:    6,
			expected: "writev$SyzLLM(0xb, 0x7ffd4b4cbca0, 0x2)",
		},
		{
			name:     "1 normal res Prov + 0 count",
			call:     "r0 = socket$inet_sctp(0x2, 0x5, 0x84)",
			count:    0,
			expected: "r0 = socket$inet_sctp(0x2, 0x5, 0x84)",
		},
		{
			name:     "1 normal res Prov + 0 count",
			call:     "r19 = socket$inet_sctp(0x2, 0x5, 0x84)",
			count:    0,
			expected: "r0 = socket$inet_sctp(0x2, 0x5, 0x84)",
		},
		{
			name:     "1 normal res Prov + 9 count",
			call:     "r3 = socket$inet_sctp(0x2, 0x5, 0x84)",
			count:    9,
			expected: "r9 = socket$inet_sctp(0x2, 0x5, 0x84)",
		},
		{
			name:     "1 multi res Prov + 0 count",
			call:     "pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})",
			count:    0,
			expected: "pipe(&(0x7f0000064000)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff})",
		},
		{
			name:     "1 multi res Prov + 9 count",
			call:     "pipe(&(0x7f0000064000)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff})",
			count:    9,
			expected: "pipe(&(0x7f0000064000)={<r9=>0xffffffffffffffff, <r10=>0xffffffffffffffff})",
		},
		{
			name:     "1 normal + 1 multi Prov",
			call:     "r19 = syscall(&(0x7f0000064000)={<r20=>0xffffffffffffffff, <r21=>0xffffffffffffffff})",
			count:    9,
			expected: "r9 = syscall(&(0x7f0000064000)={<r10=>0xffffffffffffffff, <r11=>0xffffffffffffffff})",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updateResProvider(tt.call, tt.count)
			if result != tt.expected {
				t.Errorf("updateResProvider(%q) = %v; want %v", tt.call, result, tt.expected)
			}
		})
	}
}

func TestUpdateRemainCallsResNum(t *testing.T) {
	tests := []struct {
		name        string
		calls       []string
		countAt     int
		countBefore int
		expected    []string
	}{
		{
			name:        "count 0 + offset 0",
			calls:       []string{"r0 = sendto$llc(0x2)", "r1 = openat(0x7ffff, r0)", "writev$SyzLLM(0xb, r0, r1)", "pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)", "r4 = syscall$des(0x0, r0, r1, r2, r3)"},
			countAt:     0,
			countBefore: 0,
			expected:    []string{"r0 = sendto$llc(0x2)", "r1 = openat(0x7ffff, r0)", "writev$SyzLLM(0xb, r0, r1)", "pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)", "r4 = syscall$des(0x0, r0, r1, r2, r3)"},
		},
		{
			name:        "count 0 + offset 1",
			calls:       []string{"r0 = sendto$llc(0x2)", "r1 = openat(0x7ffff, r0)", "writev$SyzLLM(0xb, r0, r1)", "pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)", "r4 = syscall$des(0x0, r0, r1, r2, r3)"},
			countAt:     1,
			countBefore: 0,
			expected:    []string{"r1 = sendto$llc(0x2)", "r2 = openat(0x7ffff, r1)", "writev$SyzLLM(0xb, r1, r2)", "pipe2$9p(&(0x7f0000000040)={<r3=>0xffffffffffffffff, <r4=>0xffffffffffffffff}, 0x0)", "r5 = syscall$des(0x0, r1, r2, r3, r4)"},
		},
		{
			name:        "count 1 + offset 0",
			calls:       []string{"r1 = sendto$llc(0x2)", "r2 = openat(0x7ffff, r1, r0)", "writev$SyzLLM(0xb, r1, r2)", "pipe2$9p(&(0x7f0000000040)={<r3=>0xffffffffffffffff, <r4=>0xffffffffffffffff}, 0x0)", "r5 = syscall$des(0x0, r1, r2, r3, r4, r0)"},
			countAt:     0,
			countBefore: 1,
			expected:    []string{"r1 = sendto$llc(0x2)", "r2 = openat(0x7ffff, r1, r0)", "writev$SyzLLM(0xb, r1, r2)", "pipe2$9p(&(0x7f0000000040)={<r3=>0xffffffffffffffff, <r4=>0xffffffffffffffff}, 0x0)", "r5 = syscall$des(0x0, r1, r2, r3, r4, r0)"},
		},
		{
			name:        "count 2 + offset 3",
			calls:       []string{"r2 = sendto$llc(0x2)", "r3 = openat(0x7ffff, r2, r0, r1)", "writev$SyzLLM(0xb, r2, r3)", "pipe2$9p(&(0x7f0000000040)={<r4=>0xffffffffffffffff, <r5=>0xffffffffffffffff}, 0x0)", "r6 = syscall$des(0x0, r2, r3, r4, r5, r0, r1)"},
			countAt:     3,
			countBefore: 2,
			expected:    []string{"r5 = sendto$llc(0x2)", "r6 = openat(0x7ffff, r5, r0, r1)", "writev$SyzLLM(0xb, r5, r6)", "pipe2$9p(&(0x7f0000000040)={<r7=>0xffffffffffffffff, <r8=>0xffffffffffffffff}, 0x0)", "r9 = syscall$des(0x0, r5, r6, r7, r8, r0, r1)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updateRemainCallsResNum(tt.calls, tt.countAt, tt.countBefore)
			if !AssertSlicesAreEqual(result, tt.expected) {
				t.Errorf("updateRemainCallsResNum(%q) = %v; want %v", tt.calls, result, tt.expected)
			}
		})
	}
}

func TestInsertCallNoResTag(t *testing.T) {
	tests := []struct {
		name       string
		calls      []string
		idx        int
		syzllmCall string
		expected   []string // same to initial calls
	}{
		{
			name: "1 Normal Prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r0 = openat$cgroup_ro(0xffffffffffffff9c, &(0x7f0000000040)='blkio.bfq.io_service_bytes_recursive\\x00', 0x275a, 0x0)", // [MASK]
				"r2 = sendto$llc(0x2)",
				"r3 = openat(0x7ffff, r2, r0, r1)",
				"writev$SyzLLM(0xb, r2, r3)",
				"pipe2$9p(&(0x7f0000000040)={<r4=>0xffffffffffffffff, <r5=>0xffffffffffffffff}, 0x0)",
				"r6 = syscall$des(0x0, r2, r3, r4, r5, r0, r1)"},
			idx:        2,
			syzllmCall: "r0 = openat$cgroup_ro(0xffffffffffffff9c, &(0x7f0000000040)='blkio.bfq.io_service_bytes_recursive\\x00', 0x275a, 0x0)",
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$cgroup_ro(0xffffffffffffff9c, &(0x7f0000000040)='blkio.bfq.io_service_bytes_recursive\\x00', 0x275a, 0x0)", // [MASK]
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)"},
		},
		{
			name: "1 Multi Prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"pipe2$9p(&(0x7f0000000040)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x0)", // [MASK]
				"r2 = sendto$llc(0x2)",
				"r3 = openat(0x7ffff, r2, r0, r1)",
				"writev$SyzLLM(0xb, r2, r3)",
				"pipe2$9p(&(0x7f0000000040)={<r4=>0xffffffffffffffff, <r5=>0xffffffffffffffff}, 0x0)",
				"r6 = syscall$des(0x0, r2, r3, r4, r5, r0, r1)"},
			idx:        2,
			syzllmCall: "pipe2$9p(&(0x7f0000000040)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x0)",
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
		},
		{
			name: "1 Normal Prov + 1 Multi Prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r0 = call$desc(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff})", // [MASK]
				"r2 = sendto$llc(0x2)",
				"r3 = openat(0x7ffff, r2, r0, r1)",
				"writev$SyzLLM(0xb, r2, r3)",
				"pipe2$9p(&(0x7f0000000040)={<r4=>0xffffffffffffffff, <r5=>0xffffffffffffffff}, 0x0)",
				"r6 = syscall$des(0x0, r2, r3, r4, r5, r0, r1)"},
			idx:        2,
			syzllmCall: "r0 = call$desc(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff})",
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = call$desc(&(0x7f0000000040)={<r3=>0xffffffffffffffff, <r4=>0xffffffffffffffff})", // [MASK]
				"r5 = sendto$llc(0x2)",
				"r6 = openat(0x7ffff, r5, r0, r1)",
				"writev$SyzLLM(0xb, r5, r6)",
				"pipe2$9p(&(0x7f0000000040)={<r7=>0xffffffffffffffff, <r8=>0xffffffffffffffff}, 0x0)",
				"r9 = syscall$des(0x0, r5, r6, r7, r8, r0, r1)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originCalls := tt.calls
			insertCallNoResTag(tt.calls, tt.idx, tt.syzllmCall)
			if !AssertSlicesAreEqual(tt.calls, tt.expected) {
				t.Errorf("insertCallNoResTag(%q) = %v; want %v", tt.calls, originCalls, tt.expected)
			}
		})
	}
}

func TestParseResource(t *testing.T) {
	tests := []struct {
		name       string
		syzllmCall string
		calls      []string
		idx        int
		expected   []string
	}{
		{
			name:       "1 res tag + no match",
			syzllmCall: "poll$SyzLLM(&(0x7f0000080000)=[{@RSTART@socket$SyzLLM(0x10, 0x3, 0x0)@REND@}], 0x1, 0x2710)",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"r4 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"poll$SyzLLM(&(0x7f0000080000)=[{r4}], 0x1, 0x2710)", // [MASK]
				"r5 = sendto$llc(0x2)",
				"r6 = openat(0x7ffff, r5, r0, r1)",
				"writev$SyzLLM(0xb, r5, r6)",
				"pipe2$9p(&(0x7f0000000040)={<r7=>0xffffffffffffffff, <r8=>0xffffffffffffffff}, 0x0)",
				"r9 = syscall$des(0x0, r5, r6, r7, r8, r0, r1)"},
		},
		{
			name:       "1 res tag + 1 match",
			syzllmCall: "poll$SyzLLM(&(0x7f0000080000)=[{@RSTART@socket$SyzLLM(0x10, 0x3, 0x0)@REND@}], 0x1, 0x2710)",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"poll$SyzLLM(&(0x7f0000080000)=[{r1}], 0x1, 0x2710)", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
		},
		{
			name:       "2 res tag + 0 match",
			syzllmCall: "poll$SyzLLM(&(0x7f0000080000)=[{@RSTART@socket$SyzLLM(0x10, 0x3, 0x0)@REND@}, {@RSTART@openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)@REND@}], 0x2, 0x0)",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"r4 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"r5 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"poll$SyzLLM(&(0x7f0000080000)=[{r4}, {r5}], 0x2, 0x0)", // [MASK]
				"r6 = sendto$llc(0x2)",
				"r7 = openat(0x7ffff, r6, r0, r1)",
				"writev$SyzLLM(0xb, r6, r7)",
				"pipe2$9p(&(0x7f0000000040)={<r8=>0xffffffffffffffff, <r9=>0xffffffffffffffff}, 0x0)",
				"r10 = syscall$des(0x0, r6, r7, r8, r9, r0, r1)"},
		},
		{
			name:       "2 res tag + 1 match",
			syzllmCall: "poll$SyzLLM(&(0x7f0000080000)=[{@RSTART@socket$SyzLLM(0x10, 0x3, 0x0)@REND@}, {@RSTART@openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)@REND@}], 0x2, 0x0)",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"pipe2$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"r4 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"poll$SyzLLM(&(0x7f0000080000)=[{r1}, {r4}], 0x2, 0x0)", // [MASK]
				"r5 = sendto$llc(0x2)",
				"r6 = openat(0x7ffff, r5, r0, r1)",
				"writev$SyzLLM(0xb, r5, r6)",
				"pipe2$9p(&(0x7f0000000040)={<r7=>0xffffffffffffffff, <r8=>0xffffffffffffffff}, 0x0)",
				"r9 = syscall$des(0x0, r5, r6, r7, r8, r0, r1)"},
		},
		{
			name:       "1 res tag + 1 pipe tag + 0 match",
			syzllmCall: "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x2, @PIPESTART@pipe2(&(0x7f0000000240)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x80800)@PIPEEND@, &(0x7f000003a000)={0x1, 0x6})",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"call$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"call$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"r4 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r4, 0x2, r5, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
				"r7 = sendto$llc(0x2)",
				"r8 = openat(0x7ffff, r7, r0, r1)",
				"writev$SyzLLM(0xb, r7, r8)",
				"pipe2$9p(&(0x7f0000000040)={<r9=>0xffffffffffffffff, <r10=>0xffffffffffffffff}, 0x0)",
				"r11 = syscall$des(0x0, r7, r8, r9, r10, r0, r1)"},
		},
		{
			name:       "1 res tag + 1 pipe tag + 1 res match",
			syzllmCall: "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x2, @PIPESTART@pipe2(&(0x7f0000000240)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x80800)@PIPEEND@, &(0x7f000003a000)={0x1, 0x6})",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"call$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"call$9p(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"pipe2(&(0x7f0000000240)={<r4=>0xffffffffffffffff, <r5=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r1, 0x2, r4, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
				"r6 = sendto$llc(0x2)",
				"r7 = openat(0x7ffff, r6, r0, r1)",
				"writev$SyzLLM(0xb, r6, r7)",
				"pipe2$9p(&(0x7f0000000040)={<r8=>0xffffffffffffffff, <r9=>0xffffffffffffffff}, 0x0)",
				"r10 = syscall$des(0x0, r6, r7, r8, r9, r0, r1)"},
		},
		{
			name:       "1 res tag + 1 pipe tag + 1 pipe match",
			syzllmCall: "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x2, @PIPESTART@pipe2(&(0x7f0000000240)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x80800)@PIPEEND@, &(0x7f000003a000)={0x1, 0x6})",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = call$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = call$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"r4 = epoll_create1$SyzLLM(0x80000)",
				"epoll_ctl$SyzLLM(r4, 0x2, r2, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
				"r5 = sendto$llc(0x2)",
				"r6 = openat(0x7ffff, r5, r0, r1)",
				"writev$SyzLLM(0xb, r5, r6)",
				"pipe2$9p(&(0x7f0000000040)={<r7=>0xffffffffffffffff, <r8=>0xffffffffffffffff}, 0x0)",
				"r9 = syscall$des(0x0, r5, r6, r7, r8, r0, r1)"},
		},
		{
			name:       "1 res tag + 1 pipe tag + 1 res match + 1 pipe match",
			syzllmCall: "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x2, @PIPESTART@pipe2(&(0x7f0000000240)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x80800)@PIPEEND@, &(0x7f000003a000)={0x1, 0x6})",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"epoll_ctl$SyzLLM(r1, 0x2, r2, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
		},
		{
			name:       "1 res tag + 1 pipe tag + 1 res match + syzllmCall has Prov",
			syzllmCall: "r0 = llm(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x2, @PIPESTART@pipe2(&(0x7f0000000240)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x80800)@PIPEEND@, &(0x7f000003a000)={0x1, 0x6})",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x0)",
				"multi(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)"},
			idx: 3,
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x0)",
				"multi(&(0x7f0000000040)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x0)",
				"pipe2(&(0x7f0000000240)={<r4=>0xffffffffffffffff, <r5=>0xffffffffffffffff}, 0x80800)",
				"r6 = llm(r1, 0x2, r4, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
				"r7 = sendto$llc(0x2)",
				"r8 = openat(0x7ffff, r7, r0, r1)",
				"writev$SyzLLM(0xb, r7, r8)",
				"pipe2$9p(&(0x7f0000000040)={<r9=>0xffffffffffffffff, <r10=>0xffffffffffffffff}, 0x0)",
				"r11 = syscall$des(0x0, r7, r8, r9, r10, r0, r1)"},
		},
		{
			name:       "1 normal call",
			syzllmCall: "bpf$MAP_CREATE(0x0, &(0x7f0000002040)=ANY=[@ANYBLOB], 0x48)",
			calls:      []string{"[MASK]"},
			idx:        0,
			expected:   []string{"bpf$MAP_CREATE(0x0, &(0x7f0000002040)=ANY=[@ANYBLOB], 0x48)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseResource(tt.calls, tt.idx, tt.syzllmCall)
			if !AssertSlicesAreEqual(result, tt.expected) {
				t.Errorf("ParseResource = \n%v\n; want \n%v", result, tt.expected)
			}
		})
	}
}

func TestParseResTag(t *testing.T) {
	tests := []struct {
		name     string
		call     string
		expected []string
	}{
		{
			name: "1 res tag",
			call: "poll$SyzLLM(&(0x7f0000080000)=[{@RSTART@socket$SyzLLM(0x10, 0x3, 0x0)@REND@}], 0x1, 0x2710)",
			expected: []string{
				"r0 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"poll$SyzLLM(&(0x7f0000080000)=[{r0}], 0x1, 0x2710)"},
		},
		{
			name: "2 res tags",
			call: "poll$SyzLLM(&(0x7f0000080000)=[{@RSTART@socket$SyzLLM(0x10, 0x3, 0x0)@REND@}, {@RSTART@openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)@REND@}], 0x2, 0x0)",
			expected: []string{
				"r0 = socket$SyzLLM(0x10, 0x3, 0x0)",
				"r1 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"poll$SyzLLM(&(0x7f0000080000)=[{r0}, {r1}], 0x2, 0x0)"},
		},
		{
			name: "1 res tag + 1 pipe tag",
			call: "epoll_ctl$SyzLLM(@RSTART@epoll_create1$SyzLLM(0x80000)@REND@, 0x2, @PIPESTART@pipe2(&(0x7f0000000240)={<r0=>0xffffffffffffffff, <r1=>0xffffffffffffffff}, 0x80800)@PIPEEND@, &(0x7f000003a000)={0x1, 0x6})",
			expected: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseResTag(tt.call)
			if !AssertSlicesAreEqual(result, tt.expected) {
				t.Errorf("parseResTag(%q) = %v; want %v", tt.call, result, tt.expected)
			}
		})
	}
}

func TestInsertSyzllmCalls(t *testing.T) {
	tests := []struct {
		name        string
		calls       []string
		insertPos   int
		syzllmCalls []string
		expected    []string
	}{
		{
			name: "no res prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"[MASK]", // [MASK]
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
			},
			insertPos: 3,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r4=>0xffffffffffffffff, <r5=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r3, 0x2, r4, &(0x7f000003a000)={0x1, 0x6})",
				"r6 = sendto$llc(0x2)",
				"r7 = openat(0x7ffff, r6, r0, r1)",
				"writev$SyzLLM(0xb, r6, r7)",
				"pipe2$9p(&(0x7f0000000040)={<r8=>0xffffffffffffffff, <r9=>0xffffffffffffffff}, 0x0)",
				"r10 = syscall$des(0x0, r6, r7, r8, r9, r0, r1)",
			},
		},
		{
			name: "1 normal res prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"[MASK]", // [MASK]
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
			},
			insertPos: 3,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"pipe2(&(0x7f0000000240)={<r3=>0xffffffffffffffff, <r4=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r1, 0x2, r3, &(0x7f000003a000)={0x1, 0x6})",
				"r5 = sendto$llc(0x2)",
				"r6 = openat(0x7ffff, r5, r0, r1)",
				"writev$SyzLLM(0xb, r5, r6)",
				"pipe2$9p(&(0x7f0000000040)={<r7=>0xffffffffffffffff, <r8=>0xffffffffffffffff}, 0x0)",
				"r9 = syscall$des(0x0, r5, r6, r7, r8, r0, r1)",
			},
		},
		{
			name: "all match: 1 normal res prov + 1 pipe res prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x80800)",
				"[MASK]", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)",
			},
			insertPos: 3,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r1, 0x2, r2, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
				"r4 = sendto$llc(0x2)",
				"r5 = openat(0x7ffff, r4, r0, r1)",
				"writev$SyzLLM(0xb, r4, r5)",
				"pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
				"r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)",
			},
		},
		{
			name: "insert to beginning",
			calls: []string{
				"[MASK]", // [MASK]
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
			},
			insertPos: 0,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
				"r3 = read$desc(&(0x7ffffff)=\"0\")",
				"r4 = dup(r3)",
				"r5 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r6 = sendto$llc(0x2)",
				"r7 = openat(0x7ffff, r6, r3, r4)",
				"writev$SyzLLM(0xb, r6, r7)",
				"pipe2$9p(&(0x7f0000000040)={<r8=>0xffffffffffffffff, <r9=>0xffffffffffffffff}, 0x0)",
				"r10 = syscall$des(0x0, r6, r7, r8, r9, r3, r4)",
			},
		},
		{
			name: "insert to end - no prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"syscall$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"[MASK]", // [MASK]
			},
			insertPos: 8,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"syscall$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"r8 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r9=>0xffffffffffffffff, <r10=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r8, 0x2, r9, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
			},
		},
		{
			name: "insert to end - 1 normal prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = epoll_create1$SyzLLM(0x80000)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"syscall$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"[MASK]", // [MASK]
			},
			insertPos: 8,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = epoll_create1$SyzLLM(0x80000)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"syscall$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"pipe2(&(0x7f0000000240)={<r8=>0xffffffffffffffff, <r9=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r3, 0x2, r8, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
			},
		},
		{
			name: "insert to end - 1 pipe prov",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"[MASK]", // [MASK]
			},
			insertPos: 8,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = sendto$llc(0x2)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"r8 = epoll_create1$SyzLLM(0x80000)",
				"epoll_ctl$SyzLLM(r8, 0x2, r5, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
			},
		},
		{
			name: "insert to end - all prov exist",
			calls: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = epoll_create1$SyzLLM(0x80000)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"[MASK]", // [MASK]
			},
			insertPos: 8,
			syzllmCalls: []string{
				"r0 = epoll_create1$SyzLLM(0x80000)",
				"pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
				"epoll_ctl$SyzLLM(r0, 0x2, r1, &(0x7f000003a000)={0x1, 0x6})",
			},
			expected: []string{
				"r0 = read$desc(&(0x7ffffff)=\"0\")",
				"r1 = dup(r0)",
				"r2 = openat$SyzLLM(0xffffffffffffff9c, &(0x7f0000008000)=nil, 0x0, 0x0)",
				"r3 = epoll_create1$SyzLLM(0x80000)",
				"r4 = openat(0x7ffff, r3, r0, r1)",
				"writev$SyzLLM(0xb, r3, r4)",
				"pipe2$9p(&(0x7f0000000040)={<r5=>0xffffffffffffffff, <r6=>0xffffffffffffffff}, 0x0)",
				"r7 = syscall$des(0x0, r3, r4, r5, r6, r0, r1)",
				"epoll_ctl$SyzLLM(r3, 0x2, r5, &(0x7f000003a000)={0x1, 0x6})", // [MASK]
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := insertSyzllmCalls(tt.calls, tt.insertPos, tt.syzllmCalls)
			if !AssertSlicesAreEqual(result, tt.expected) {
				t.Errorf("removeExistingResProvider = \n%v\n; want \n%v", result, tt.expected)
			}
		})
	}
}

func TestExtractCallName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "r0 = read$desc(&(0x7ffffff)=\"0\")",
			expected: "read",
		},
		{
			input:    "r1 = epoll_create1$SyzLLM(0x80000)",
			expected: "epoll_create1",
		},
		{
			input:    "pipe2(&(0x7f0000000240)={<r2=>0xffffffffffffffff, <r3=>0xffffffffffffffff}, 0x80800)",
			expected: "pipe2",
		},
		{
			input:    "epoll_ctl$SyzLLM(r1, 0x2, r2, &(0x7f000003a000)={0x1, 0x6})",
			expected: "epoll_ctl",
		},
		{
			input:    "r4 = sendto$llc(0x2)",
			expected: "sendto",
		},
		{
			input:    "r5 = openat(0x7ffff, r4, r0, r1)",
			expected: "openat",
		},
		{
			input:    "writev$SyzLLM(0xb, r4, r5)",
			expected: "writev",
		},
		{
			input:    "pipe2$9p(&(0x7f0000000040)={<r6=>0xffffffffffffffff, <r7=>0xffffffffffffffff}, 0x0)",
			expected: "pipe2",
		},
		{
			input:    "r8 = syscall$des(0x0, r4, r5, r6, r7, r0, r1)",
			expected: "syscall",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractCallName(tt.input)
			if result != tt.expected {
				t.Errorf("extractCallName(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractFirstResProvider(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "r0 = epoll_create1$SyzLLM(0x80000)",
			expected: "r0",
		},
		{
			input:    "pipe2(&(0x7f0000000240)={<r1=>0xffffffffffffffff, <r2=>0xffffffffffffffff}, 0x80800)",
			expected: "r1",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "no_resource_tag()",
			expected: "",
		},
		{
			input:    "r3 = invalid<no_tag>",
			expected: "r3",
		},
		{
			input:    "<r4=>0x0, <r5=>0x0}",
			expected: "r4",
		},
		{
			input:    "r6=malformed",
			expected: "r6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractFirstResProvider(tt.input)
			if result != tt.expected {
				t.Errorf("extractFirstResourceTag(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
