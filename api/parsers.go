package api

import (
	"fmt"
	"os/exec"
	"strings"
)

type ParserInterface interface {
	Parse(string) (*map[string]interface{}, error)
}

type Parser struct {
	path         string
	container_id string
}

func (self *Parser) GetPath() string {
	return self.path
}

func (self *Parser) Parse(container_id *string) (map[string]interface{}, error) {

	var file = fmt.Sprintf(self.path, (*container_id))
	var cmd = exec.Command("cat", file)
	cat, err := cmd.Output()
	HandleErrorAndExitWithMsg(err, fmt.Sprintf("File does not exist - %s", file))

	var data = make(map[string]interface{})
	for _, line := range strings.Split(string(cat), "\n") {
		var key_val = strings.Split(line, " ")
		var k = strings.TrimSpace(strings.Join(key_val[:1], ""))
		var v = strings.TrimSpace(strings.Join(key_val[1:], ""))
		data[k] = v
	}

	return data, nil
}

type CpuParser struct {
	*Parser
}

func NewCpuParser() *CpuParser {
	return &CpuParser{
		Parser: &Parser{
			path: "/sys/fs/cgroup/cpuacct/docker/%s/cpuacct.stat",
		},
	}
}

type MemoryParser struct {
	*Parser
}

func NewMemoryParser() *MemoryParser {
	return &MemoryParser{
		Parser: &Parser{
			path: "/sys/fs/cgroup/memory/docker/%s/memory.stat",
		},
	}
}
