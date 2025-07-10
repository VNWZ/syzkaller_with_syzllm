package prog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/syzkaller/pkg/log"
	syzllm_pkg2 "github.com/google/syzkaller/pkg/syzllm_pkg"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func newSyzllm(prog *Prog, insertPosition int, choiceTable *ChoiceTable) syzllm {
	s := &baseSyzllm{
		program:        prog,
		insertPosition: insertPosition,
		choiceTable:    choiceTable,

		result:           prog,
		callListWithMask: nil,
		syzllmCall:       "",
	}
	return s
}

type syzllm interface {
	prepareForRequest()
	request()
	processSyzLLMCall()
	insert() *Prog
}

type baseSyzllm struct {
	program        *Prog
	insertPosition int
	choiceTable    *ChoiceTable

	result *Prog

	callListWithMask []string
	syzllmCall       string
}

func (s *baseSyzllm) insert() *Prog {
	s.prepareForRequest()
	s.request()
	s.processSyzLLMCall()
	return s.result
}

func (s *baseSyzllm) prepareForRequest() {
	callList := s.parseProgToSlice()

	var err error
	s.callListWithMask, err = syzllm_pkg2.InsertItem(callList, "[MASK]", s.insertPosition)
	if err != nil {
		panic(err)
	}
}

func (s *baseSyzllm) request() {
	serverHostInDocker := os.Getenv("SERVER_HOST")
	serverPortInDocker := os.Getenv("SERVER_PORT")

	serverInfo, err := syzllm_pkg2.GetServerInfo("test")
	if err != nil {
		panic(err)
	}
	if serverHostInDocker != "" && serverPortInDocker != "" {
		serverInfo.Host = serverHostInDocker
		serverInfo.Port = serverPortInDocker
	}
	url := fmt.Sprintf("http://%s:%s", serverInfo.Host, serverInfo.Port)

	jsonData, err := json.Marshal(syzllm_pkg2.SyscallRequestData{Syscalls: s.callListWithMask})
	if err != nil {
		panic(err)
	}

	var resp *http.Response
	responseChan, errorChan := syzllm_pkg2.SendPostRequestAsync(url, jsonData)
	select {
	case respC := <-responseChan:
		if respC != nil {
			resp = respC
		}
	case err := <-errorChan:
		if err != nil {
			fmt.Println("Response Error:", err)
			s.result = s.program
			return
		}
	}

	responseByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Response Reading Error:", err)
		s.result = s.program
		return
	}

	syzLLMResponse := syzllm_pkg2.SyzLLMResponse{State: -1, Syscall: ""}
	err = json.Unmarshal(responseByte, &syzLLMResponse)
	if err != nil {
		fmt.Println("Unmarshal Response Error:", err)
		s.result = s.program
		return
	}
	if syzLLMResponse.State != 0 {
		s.result = s.program
		return
	}

	s.syzllmCall = syzLLMResponse.Syscall
}

func (s *baseSyzllm) processSyzLLMCall() {
	s.syzllmCall = syzllm_pkg2.ProcessDescriptor(s.syzllmCall)
	resultCalls := syzllm_pkg2.ParseResource(s.callListWithMask, s.insertPosition, s.syzllmCall)

	resultText := ""
	for _, call := range resultCalls {
		if len(resultText) > 0 {
			resultText += "\n"
		}
		resultText += call
	}
	callBytes := []byte(resultText)
	prog, err := s.program.Target.Deserialize(callBytes, NonStrict)
	if err != nil {
		s.result = s.program
		return
	}

	s.result = prog
}

func (s *baseSyzllm) parseProgToSlice() []string {
	syscallBytes := s.program.Serialize()
	syscallList := strings.Split(string(syscallBytes[:]), "\n")
	if syscallList[len(syscallList)-1] == "" {
		syscallList = syscallList[:len(syscallList)-1]
	}

	return syscallList
}

func VerifyCallsFrom(path string, target *Target) {
	dummyProg := &Prog{
		Target:   target,
		Calls:    make([]*Call, 0),
		Comments: make([]string, 0),
		isUnsafe: false,
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// for debug attaching
	time.Sleep(20 * time.Second)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}
		line = strings.TrimSuffix(line, "\n")

		line = syzllm_pkg2.ProcessDescriptor(line)

		newCall := line
		maskedSyscallList := make([]string, 1)
		maskedSyscallList[0] = "[MASK]"
		calls := syzllm_pkg2.ParseResource(maskedSyscallList, 0, newCall)
		newSyscallSequence := ""
		for _, call := range calls {
			if len(newSyscallSequence) > 0 {
				newSyscallSequence += "\n"
			}
			newSyscallSequence += call
		}
		newSyscallBytes := []byte(newSyscallSequence)
		_, err = dummyProg.Target.Deserialize(newSyscallBytes, NonStrict)
		if err != nil {
			log.Logf(0, "invlid calls: "+newCall, err)
		}
	}
	log.Fatalf("Verify done")
}
