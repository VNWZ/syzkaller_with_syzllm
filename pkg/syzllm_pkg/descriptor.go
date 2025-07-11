package syzllm_pkg

import (
	"github.com/google/syzkaller/pkg/log"
	"regexp"
	"strings"
)

func ExtractCallNameWithoutDescriptor(call string) string {
	CallNamePattern := regexp.MustCompile(`([a-zA-Z0-9_]+)\$`)
	match := CallNamePattern.FindStringSubmatch(call)

	if len(match) > 1 {
		return match[1]
	} else {
		log.Fatalf("wrong resource call")
	}
	return ""
}

func ProcessDescriptor(line string) string {
	if !strings.Contains(line, "$SyzLLM") {
		return line
	}

	callsWithoutDescriptor := []string{"pipe"}
	for _, c := range callsWithoutDescriptor {
		if strings.HasPrefix(line, c) {
			return line
		}
	}

	callsWithDescriptor := map[string]string{"socketpair": "unix", "socket": "unix", "connect": "unix", "getsockname": "unix", "openat": "damon_target_ids", "fcntl": "setflags", "accept4": "unix", "read": "FUSE", "quotactl": "Q_QUOTAON", "getpeername": "llc", "sendmmsg": "unix", "getsockopt": "kcm_KCM_RECV_DISABLE", "bind": "unix", "sendto": "llc", "setsockopt": "kcm_KCM_RECV_DISABLE", "write": "damon_target_ids", "ioctl": "FITRIM", "mmap": "IORING_OFF_SQ_RING", "recvmsg": "unix", "sendmsg": "unix", "epoll_ctl": "EPOLL_CTL_ADD", "accept": "unix", "prctl": "PR_SET_PDEATHSIG", "recvfrom": "unix", "syz_open_dev": "ttys"}

	callName := ExtractCallNameWithoutDescriptor(line)
	descriptor, exists := callsWithDescriptor[callName]
	if exists {
		line = strings.Replace(line, "$SyzLLM", "$"+descriptor, 1)
	} else {
		line = strings.Replace(line, "$SyzLLM", "", 1)
	}

	line = ProcessNestedDescriptor(line, callsWithDescriptor)

	return line
}

func ProcessNestedDescriptor(data string, replacements map[string]string) string {
	pattern := `@RSTART@(.*?)\$SyzLLM`
	re := regexp.MustCompile(pattern)

	replacer := func(match string) string {
		syscallName := strings.Split(match, "$SyzLLM")[0]
		syscallName = strings.TrimPrefix(syscallName, "@RSTART@")

		descriptor, ok := replacements[syscallName]
		if ok {
			return "@RSTART@" + syscallName + "$" + descriptor
		} else {
			return "@RSTART@" + syscallName
		}
	}

	// Replace all occurrences using the replacer function
	result := re.ReplaceAllStringFunc(data, replacer)
	return result
}
