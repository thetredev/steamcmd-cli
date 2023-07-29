package daemon

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/creack/pty"
)

const serverConsoleInputDelay time.Duration = 250 * time.Millisecond

var logger log.Logger = *log.Default()

var serverCommand *exec.Cmd = nil

var serverConsole *os.File = nil
var serverConsoleInput chan string
var serverConsoleOutput []string

func init() {
	serverConsoleInput = make(chan string)
}

func isSRCDS() bool {
	// check whether Config.ServerHome contains the executable 'srcds_linux'
	return false
}

func isCSGO() bool {
	// isSRCDS + check whether Config.ServerHome contains the folder 'csgo'
	return false
}

func serverConsoleInputListener(commandChannel <-chan string) {
	for {
		select {
		case command := <-commandChannel:
			serverConsole.WriteString(fmt.Sprintf("%s\n", command))

			if command == "quit" {
				serverCommand.Wait()
				serverCommand = nil
				serverConsole = nil
				serverConsoleOutput = make([]string, 0)
				return
			}
		}
	}
}

func serverConsoleOutputListener() {
	scanner := bufio.NewScanner(serverConsole)

	for scanner.Scan() {
		serverConsoleOutput = append(serverConsoleOutput, scanner.Text())
	}
}

func SendServerLogs(receiver *net.UDPAddr) {
	logger.Printf("Received request to send game server logs to %v\n", receiver)

	if serverConsoleOutput != nil {
		logger.Printf("Sending game server logs to %v...\n", receiver)

		bytes := 0

		for _, line := range serverConsoleOutput {
			SendSocketResponseMessage(receiver, line)
			bytes += len(line)
		}

		logger.Printf("Sent %d bytes (%d lines) of game server logs to %v\n", bytes, len(serverConsoleOutput), receiver)
	} else {
		logger.Println("Ignoring: Nothing to send.")
	}
}

func sendConsoleReplies(receiver *net.UDPAddr, command string) {
	logger.Printf("Sending game server console replies for command '%s' to %v...\n", command, receiver)

	// ensure the console replies are printed to as expected
	time.Sleep(serverConsoleInputDelay)

	// get all console replies in reverse order to save CPU cycles and memory
	reversedLines := []string{}

	for i := len(serverConsoleOutput) - 1; i > -1; i-- {
		line := serverConsoleOutput[i]

		if line == command {
			break
		}

		reversedLines = append(reversedLines, line)
	}

	// send all console replies to receiver in correct order
	for i := len(reversedLines) - 1; i > -1; i-- {
		SendSocketResponseMessage(receiver, reversedLines[i])
	}
}

func SendConsoleCommand(receiver *net.UDPAddr, command string) {
	logger.Printf("Received server command '%s' from %v\n", command, receiver)

	if serverCommand != nil && serverConsole != nil {
		logger.Printf("Sending server command '%s' from %v to game server console...\n", command, receiver)

		serverConsoleInput <- command
		sendConsoleReplies(receiver, command)
	}
}

func StartServer(receiver *net.UDPAddr) {
	logger.Printf("Received request to start the game server from %v\n", receiver)

	if serverCommand != nil {
		message := fmt.Sprintf("Server already running (PID: %d)", serverCommand.Process.Pid)

		logger.Printf("Ignoring: %s\n", message)
		SendSocketResponseMessage(receiver, message)
		return
	}

	if _, err := os.Stat(Config.ServerHome); os.IsNotExist(err) {
		log.Fatal("STEAMCMD_SERVER_HOME is set to a nonexistent path")
	}

	serverCommand = exec.Command(
		"bash", "-c",
		"./srcds_linux -console -game cstrike +ip 0.0.0.0 -port 27019 +maxplayers 12 +map de_dust2 -tickrate 128 -threads 3 -nodev",
	)

	serverCommand.Env = os.Environ()
	serverCommand.Env = append(serverCommand.Env,
		fmt.Sprintf("LD_LIBRARY_PATH=.:./bin:%s", os.Getenv("LD_LIBRARY_PATH")),
		"RESTART=no",
	)

	serverCommand.Dir = Config.ServerHome

	logger.Println("Starting game server...")

	var err error
	serverConsole, err = pty.Start(serverCommand)

	if err != nil {
		log.Fatal(err)
	}

	go serverConsoleInputListener(serverConsoleInput)

	serverConsoleOutput = []string{}
	go serverConsoleOutputListener()
}

func StopServer(receiver *net.UDPAddr) {
	logger.Printf("Received request to stop the game server from %v\n", receiver)
	SendConsoleCommand(receiver, "quit")
}

func UpdateServer(receiver *net.UDPAddr) {
	logger.Printf("Received request to update the game server from %v\n", receiver)

	if serverCommand != nil && serverConsole != nil {
		SendSocketResponseMessage(receiver, "Server is running, cannot update. Ignoring...")
		return
	}

	if len(Config.Application) == 0 {
		log.Fatal("STEAMCMD_SH not set")
	}

	if _, err := os.Stat(Config.Application); os.IsNotExist(err) {
		log.Fatal(err)
	}

	if Config.ServerAppId <= 0 {
		log.Fatal("STEAMCMD_SERVER_APPID not set")
	}

	if len(Config.ServerHome) == 0 {
		log.Fatal("STEAMCMD_SERVER_HOME not set")
	}

	updateCommand := exec.Command(
		Config.Application,
		"+force_install_dir", Config.ServerHome,
		"+login", "anonymous",
		Config.ServerAppConfig,
		"+app_update", strconv.Itoa(Config.ServerAppId), "validate",
		"+quit",
	)

	updateCommand.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := updateCommand.StdoutPipe()
	updateCommand.Stderr = updateCommand.Stdout

	if err != nil {
		log.Fatal(err)
	}

	logger.Println("Updating the game server...")

	if err = updateCommand.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			line := scanner.Text()
			SendSocketResponseMessage(receiver, line)
		}
	}()

	updateCommand.Wait()
}
