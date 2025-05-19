package prog

import (
	"encoding/json"
	"fmt"
	syzllm_pkg2 "github.com/google/syzkaller/pkg/syzllm_pkg"
	"io"
	"net/http"
	"strings"
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
	serverInfo, err := syzllm_pkg2.GetServerInfo("test")
	if err != nil {
		panic(err)
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
	resultCalls := syzllm_pkg2.ParseResource(s.syzllmCall, s.callListWithMask, s.insertPosition)

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
