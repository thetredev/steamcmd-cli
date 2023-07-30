package daemon

import (
	"bufio"
	"fmt"
	"os"
)

type ServerConsole struct {
	Pty    *os.File
	Input  chan string
	Output []string
}

func NewServerConsole() *ServerConsole {
	return &ServerConsole{
		Pty:    nil,
		Input:  make(chan string),
		Output: make([]string, 0),
	}
}

func (console *ServerConsole) Delete() {
	console.Pty = nil
	console.Output = make([]string, 0)
}

func (console *ServerConsole) AppendOutputLine(line string) {
	console.Output = append(console.Output, line)
}

func (console *ServerConsole) ListenForInput() {
	for {
		select {
		case input := <-console.Input:
			console.Pty.WriteString(fmt.Sprintf("%s\n", input))

			if input == "quit" {
				return
			}
		}
	}
}

func (console *ServerConsole) ListenForOutput() {
	console.Output = []string{}
	scanner := bufio.NewScanner(console.Pty)

	for scanner.Scan() {
		console.AppendOutputLine(scanner.Text())
	}
}

func (console *ServerConsole) SendCommandReplies(socket *Socket, command string) {
	// get all console replies in reverse order to save CPU cycles and memory
	reversedLines := []string{}

	for i := len(console.Output) - 1; i > -1; i-- {
		line := console.Output[i]

		if line == command {
			break
		}

		reversedLines = append(reversedLines, line)
	}

	// send all console replies to receiver in correct order
	for i := len(reversedLines) - 1; i > -1; i-- {
		socket.Input <- reversedLines[i]
	}
}

func (console *ServerConsole) SendLogs(socket *Socket) int {
	bytes := 0

	for _, line := range console.Output {
		socket.Input <- line
		bytes += len(line)
	}

	return bytes
}
