package daemon

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

const serverConsoleInputDelay time.Duration = 250 * time.Millisecond

type Server struct {
	Logger  *log.Logger
	Command *exec.Cmd
	Console *ServerConsole
}

func NewServer() *Server {
	return &Server{
		Logger:  log.Default(),
		Command: nil,
		Console: NewServerConsole(),
	}
}

func (server *Server) Delete() {
	if server.Console != nil {
		server.Console.Delete()
	}

	server.Command = nil
}

func (server *Server) Update(socket *Socket) error {
	server.Logger.Println("Received request to update the game server")

	if server.IsRunning() {
		message := "Server is currently running, cannot update."

		socket.SendMessage(message)
		server.Logger.Println("Ignoring:", message)
		return nil
	}

	server.Logger.Println("Updating the game server...")

	if len(Config.Application) == 0 {
		return errors.New("STEAMCMD_SH not set")
	}

	if _, err := os.Stat(Config.Application); os.IsNotExist(err) {
		return err
	}

	if Config.ServerAppId <= 0 {
		return errors.New("STEAMCMD_SERVER_APPID not set")
	}

	if len(Config.ServerHome) == 0 {
		return errors.New("STEAMCMD_SERVER_HOME not set")
	}

	updateCommand := exec.Command(
		"bash", "-c",
		fmt.Sprintf(
			"%s +force_install_dir %s +login anonymous %s +app_update %d validate +quit",
			Config.Application,
			Config.ServerHome,
			Config.ServerAppConfig,
			Config.ServerAppId,
		),
	)

	ptyFile, err := pty.Start(updateCommand)

	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(ptyFile)

		for scanner.Scan() {
			socket.SendMessage(scanner.Text())
		}
	}()

	updateCommand.Wait()
	return nil
}

func gameString(server *Server) string {
	if server.IsSRCDS() && Config.ServerGame == "css" {
		return "cstrike"
	}

	return Config.ServerGame
}

func maxplayersString(server *Server) string {
	if server.IsCSGO() {
		return "-maxplayers_override"
	}

	return "+maxplayers"
}

func (server *Server) Start(socket *Socket) error {
	server.Logger.Println("Received request to start the game server")

	if server.IsRunning() {
		message := fmt.Sprintf("Server already running (PID: %d)", server.Command.Process.Pid)

		socket.SendMessage(message)
		server.Logger.Println("Ignoring:", message)
	}

	if _, err := os.Stat(Config.ServerHome); os.IsNotExist(err) {
		return errors.New("STEAMCMD_SERVER_HOME is set to a nonexistent path")
	}

	if len(Config.ServerGame) == 0 {
		return errors.New("STEAMCMD_SERVER_GAME not set")
	}

	if Config.ServerMaxPlayers == 0 {
		return errors.New("STEAMCMD_SERVER_MAXPLAYERS not set")
	}

	if len(Config.ServerMap) == 0 {
		return errors.New("STEAMCMD_SERVER_MAP not set")
	}

	server.Logger.Println("Starting game server...")

	server.Command = exec.Command(
		"bash", "-c",
		fmt.Sprintf(
			"./srcds_linux -console -game %s +ip 0.0.0.0 -port %d %s %d +map %s -tickrate %d -threads %d -nodev",
			gameString(server),
			Config.ServerPort,
			maxplayersString(server),
			Config.ServerMaxPlayers,
			Config.ServerMap,
			Config.ServerTickrate,
			Config.ServerThreads,
		),
	)

	server.Command.Env = os.Environ()
	server.Command.Env = append(server.Command.Env,
		fmt.Sprintf("LD_LIBRARY_PATH=.:./bin:%s", os.Getenv("LD_LIBRARY_PATH")),
		"RESTART=no",
	)

	server.Command.Dir = Config.ServerHome

	var err error
	server.Console.Pty, err = pty.Start(server.Command)

	if err != nil {
		return err
	}

	go server.Console.ListenForInput()
	go server.Console.ListenForOutput()

	return nil
}

func (server *Server) Stop(socket *Socket) {
	server.Logger.Println("Received request to stop the game server")
	server.DispatchConsoleCommand(socket, "quit")
}

func (server *Server) IsRunning() bool {
	return server.Command != nil && server.Console != nil
}

func pathExists(path string) bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s", Config.ServerHome, path))
	return !os.IsNotExist(err)
}

func (server *Server) IsSRCDS() bool {
	return pathExists("srcds_linux")
}

func (server *Server) IsCSGO() bool {
	return server.IsSRCDS() && pathExists("csgo")
}

func (server *Server) SendLogs(socket *Socket) {
	server.Logger.Println("Received request to send game server logs")

	if server.IsRunning() {
		server.Logger.Println("Sending game server logs")

		bytes := server.Console.SendLogs(socket)
		server.Logger.Printf("Sent %d bytes (%d lines) of game server logs", bytes, len(server.Console.Output))
	} else {
		server.Logger.Println("Ignoring: Nothing to send.")
	}
}

func (server *Server) DispatchConsoleCommand(socket *Socket, command string) {
	server.Logger.Printf("Received server command '%s'\n", command)

	if server.IsRunning() {
		server.Logger.Printf("Sending server command '%s' to game server console...\n", command)
		server.Console.Input <- command

		if command == "quit" {
			server.Command.Wait()
			server.Delete()
		} else {
			// ensure the console replies are printed to as expected
			time.Sleep(serverConsoleInputDelay)

			server.Logger.Printf("Sending game server console replies for command '%s'...\n", command)
			server.Console.SendCommandReplies(socket, command)
		}
	}
}
