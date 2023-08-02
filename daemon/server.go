package daemon

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	"golang.org/x/exp/slices"
)

const SERVER_CONSOLE_INPUT_DELAY time.Duration = 250 * time.Millisecond

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
		server.Console = nil
	}

	server.Command = nil
}

func gameServerMod() string {
	if Config.ServerAppId != 90 || Config.ServerMod == "valve" {
		return ""
	}

	return fmt.Sprintf("+app_set_config 90 mod %s", Config.ServerMod)
}

func (server *Server) Update(socket *Socket) (bool, error) {
	server.Logger.Println("Received request to update the game server")

	if server.IsRunning() {
		message := "Server is currently running, cannot update."

		socket.SendMessage(message)
		server.Logger.Println("Ignoring:", message)
		return true, nil
	}

	server.Logger.Println("Updating the game server...")

	if len(Config.Application) == 0 {
		return false, errors.New("STEAMCMD_SH not set")
	}

	if _, err := os.Stat(Config.Application); os.IsNotExist(err) {
		return false, err
	}

	if Config.ServerAppId <= 0 {
		return false, errors.New("STEAMCMD_SERVER_APPID not set")
	}

	if len(Config.ServerHome) == 0 {
		return false, errors.New("STEAMCMD_SERVER_HOME not set")
	}

	updateCommand := exec.Command(
		"bash", Config.Application,
		"+force_install_dir", Config.ServerHome,
		"+login", "anonymous",
		gameServerMod(),
		"+app_update", fmt.Sprintf("%d", Config.ServerAppId), "validate",
		"+quit",
	)

	updater := NewServerUpdater()

	var err error
	updater.Pty, err = pty.Start(updateCommand)

	if err != nil {
		return false, err
	}

	go func() {
		scanner := bufio.NewScanner(updater.Pty)

		for scanner.Scan() {
			line := scanner.Text()

			updater.AppendOutputLine(line)
			socket.SendMessage(line)
		}
	}()

	updateCommand.Wait()
	success := updater.IsSuccessful()

	updater.Delete()
	return success, nil
}

func gameServerString(server *Server) string {
	if server.IsSRCDS() {
		return "srcds_linux"
	}

	return "hlds_linux"
}

func gameString(server *Server) (string, error) {
	if Config.ServerMod == "valve" {
		return "valve", nil
	}

	entries, _ := os.ReadDir(Config.ServerHome)
	invalidEntries := []string{
		"bin",
		"hl2",
		"linux",
		"linux32",
		"linux64",
		"platform",
		"steamapps",
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()

		if slices.Contains(invalidEntries, dirName) {
			continue
		}

		if strings.Contains(dirName, "_") || strings.Contains(dirName, ".pak") {
			continue
		}

		return dirName, nil
	}

	return "", errors.New("game not found, try 'update' first")
}

func maxplayersString(server *Server) string {
	if server.IsCSGO() {
		return "-maxplayers_override"
	}

	return "+maxplayers"
}

func fpsMaxString(server *Server) string {
	if Config.ServerFpsMax == 0 {
		return ""
	}

	return fmt.Sprintf("+fps_max %d", Config.ServerFpsMax)
}

func (server *Server) Start(socket *Socket) error {
	server.Logger.Println("Received request to start the game server")

	if server.IsRunning() {
		message := fmt.Sprintf("Server already running (PID: %d)", server.Command.Process.Pid)

		socket.SendMessage(message)
		server.Logger.Println("Ignoring:", message)
	}

	var err error
	if _, err = os.Stat(Config.ServerHome); os.IsNotExist(err) {
		return errors.New("STEAMCMD_SERVER_HOME is set to a nonexistent path")
	}

	if Config.ServerMaxPlayers == 0 {
		return errors.New("STEAMCMD_SERVER_MAXPLAYERS not set")
	}

	if len(Config.ServerMap) == 0 {
		return errors.New("STEAMCMD_SERVER_MAP not set")
	}

	var game string
	game, err = gameString(server)

	if err != nil {
		return err
	}

	server.Logger.Println("Starting game server...")

	server.Command = exec.Command(
		"bash", "-c",
		fmt.Sprintf(
			"./%s -console -game %s +ip 0.0.0.0 -port %d %s %d +map %s %s -tickrate %d -threads %d -nodev",
			gameServerString(server),
			game,
			Config.ServerPort,
			maxplayersString(server),
			Config.ServerMaxPlayers,
			Config.ServerMap,
			fpsMaxString(server),
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
	server.Console.Pty, err = pty.Start(server.Command)

	if err != nil {
		return err
	}

	go server.Console.ListenForInput()
	go server.Console.ListenForOutput()

	if server.IsCSGO() {
		server.EnableTickrate()
	}

	socket.SendMessage("Game server started. You can now view its logs.")
	return nil
}

func (server *Server) Stop(socket *Socket) {
	server.Logger.Println("Received request to stop the game server")

	if server.IsRunning() {
		server.DispatchConsoleCommand(socket, "quit")

		time.Sleep(TCP_CONGESTION_PREVENTION_DELAY)
		socket.SendMessage("Game server stopped.")
	} else {
		message := "Game server not running. Nothing to stop."

		socket.SendMessage(message)
		server.Logger.Printf("Ignoring: %s\n", message)
	}
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
		message := "Game server not running. Nothing to send."

		socket.SendMessage(message)
		server.Logger.Printf("Ignoring: %s\n", message)
	}
}

func (server *Server) DispatchConsoleCommand(socket *Socket, command string) {
	server.Logger.Printf("Received server command '%s'\n", command)

	if server.IsRunning() {
		server.Logger.Printf("Sending server command '%s' to game server console...\n", command)
		server.Console.Input <- command

		message := fmt.Sprintf("Command '%s' dispatched to game server console.", command)

		if command == "quit" {
			server.Command.Wait()
			server.Delete()

			socket.SendMessage(message)
		} else {
			socket.SendMessage(fmt.Sprintln(message))

			// ensure the console replies are printed to as expected
			time.Sleep(SERVER_CONSOLE_INPUT_DELAY)

			server.Logger.Printf("Sending game server console replies for command '%s'...\n", command)
			server.Console.SendCommandReplies(socket, command)
		}
	}
}

func (server *Server) EnableTickrate() {
	var tickrateCvars = []string{
		"sv_minupdaterate",
		"sv_mincmdrate",
		"sv_minrate",
		"sv_maxrate",
	}

	for _, tickrateCvar := range tickrateCvars {
		server.Console.Input <- fmt.Sprintf("%s %d", tickrateCvar, Config.ServerTickrate)
	}
}
