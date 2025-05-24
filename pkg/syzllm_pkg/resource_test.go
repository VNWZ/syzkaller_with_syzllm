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
