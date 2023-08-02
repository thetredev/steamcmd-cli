package daemon

import (
	"fmt"
	"os"
	"strings"
)

type ServerUpdater struct {
	Pty    *os.File
	Output []string
}

func NewServerUpdater() *ServerUpdater {
	return &ServerUpdater{
		Pty:    nil,
		Output: make([]string, 0),
	}
}

func (updater *ServerUpdater) Delete() {
	updater.Pty = nil
	updater.Output = make([]string, 0)
}

func (updater *ServerUpdater) AppendOutputLine(line string) {
	updater.Output = append(updater.Output, line)
}

func (updater *ServerUpdater) IsSuccessful() bool {
	var successMessage = fmt.Sprintf("Success! App '%d' fully installed.", Config.ServerAppId)

	for i := len(updater.Output) - 1; i > -1; i-- {
		if strings.Contains(updater.Output[i], successMessage) {
			return true
		}
	}

	return false
}
