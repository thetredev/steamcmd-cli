# steamcmd-cli

A custom SteamCMD client and game server manager implementation written in Go.

Only Linux is supported. On macOS and Windows you can build and use container images.

## Introduction

The tool was designed with a client/server architecture in mind, though the "server" part is called *daemon* and the "client" part is called *server*, in this case. The reason these terms were chosen is simple: The *daemon* listens for incoming requests to manage the game server, and the client issues *server* commands.

Hence, the tool is split into the following two main subcommands:
```
$ steamcmd-cli daemon
$ steamcmd-cli server
```

You'll get the hang of it once you've used the tool a bit.

## TCP socket communication

Communication between the *deamon* and the *server* is accomplished via a TCP socket over SSL. This requires valid certificates and certificate keys to be setup beforehand.

# Table of Contents

- [CLI commands overview](#cli-commands-overview)
- [CLI demo](#cli-demo)
- [Installation](#installation)
- [Default configuration](#default-configuration)
- [Certificate configuration](#certificate-configuration)
- [Using container images](#using-container-images)
- [Shell auto-completions](#shell-auto-completions)

# CLI commands overview

See [./docs/commands.md](./docs/commands.md)


## CLI demo

Let the CLI speak for itself (certificates are omitted for brevity, see [Default configuration](#default-configuration) and [Certificate Configuration](#certificate-configuration) for details):
```
# [shell DAEMON]

# export some environment variables
$ export STEAMCMD_SH=/usr/games/steamcmd
$ export STEAMCMD_SERVER_HOME=~/servers/css
$ export STEAMCMD_SERVER_APPID=232330
$ export STEAMCMD_SERVER_MAP=de_dust2
$ export STEAMCMD_SERVER_MAXPLAYERS=16

# fps_max is not necessary, but you can set it as well:
$ export STEAMCMD_SERVER_FPSMAX=300

# start the daemon on default port 65000
$ steamcmd-cli daemon
2023/08/02 19:15:28 Listening for incoming requests on port 65000/TCP...
```

Start up another shell and use the client:
```
# [shell SERVER]

# set the daemon IP
$ export STEAMCMD_CLI_SOCKET_IP=127.0.0.1

# optionally set the port as well (this is only needed if the same was passed to the daemon):
$ export STEAMCMD_CLI_SOCKET_PORT=65000

# update the server
$ steamcmd-cli server update
WARNING: setlocale('en_US.UTF-8') failed, using locale: 'C'. International characters may not work.
Redirecting stderr to '/home/Steam/logs/stderr.txt'
[  0%] Checking for available updates...
[----] Verifying installation...
Steam Console Client (c) Valve Corporation - version 1690585855
-- type 'quit' to exit --
Loading Steam API...dlmopen steamservice.so failed: steamservice.so: cannot open shared object file: No such file or directory
OK

Connecting anonymously to Steam Public...OK
Waiting for client config...OK
Waiting for user info...OK
 Update state (0x3) reconfiguring, progress: 0.00 (0 / 0)
 Update state (0x3) reconfiguring, progress: 0.00 (0 / 0)
......
```

Check daemon shell:
```
# [shell DAEMON]

2023/08/02 19:18:12 Received request to update the game server
2023/08/02 19:18:12 Updating the game server...
```

When the update is finished, you can start the server and view its logs:
```
# [shell SERVER]

$ steamcmd-cli server start
Game server started. You can now view its logs.

$ steamcmd-cli logs
Cannot read termcap database;
using dumb terminal settings.
Using Breakpad minidump system. Version: 6630498 AppID: 232330
Setting breakpad minidump AppID = 232330
Using breakpad crash handler
Loaded 1335 VPK file hashes from /steamcmd/server/cstrike/cstrike_pak.vpk for pure server operation.
Loaded 1335 VPK file hashes from /steamcmd/server/cstrike/cstrike_pak.vpk for pure server operation.
Loaded 1218 VPK file hashes from /steamcmd/server/hl2/hl2_textures.vpk for pure server operation.
Loaded 574 VPK file hashes from /steamcmd/server/hl2/hl2_sound_vo_english.vpk for pure server operation.
Loaded 383 VPK file hashes from /steamcmd/server/hl2/hl2_sound_misc.vpk for pure server operation.
Loaded 446 VPK file hashes from /steamcmd/server/hl2/hl2_misc.vpk for pure server operation.
Loaded 5 VPK file hashes from /steamcmd/server/platform/platform_misc.vpk for pure server operation.
server_srv.so loaded for "Counter-Strike: Source"
maxplayers set to 16
Looking up breakpad interfaces from steamclient
Calling BreakpadMiniDumpSystemInit
08/02 19:21:02 Init: Installing breakpad exception handler for appid(232330)/version(6630498)/tid(110)
Unknown command "r_decal_cullsize"
maxplayers set to 16
ConVarRef dev_loadtime_map_start doesn't point to an existing ConVar
Network: IP 0.0.0.0, mode MP, dedicated Yes, ports 27015 SV / 27005 CL
Initializing Steam libraries for secure Internet server
[S_API] SteamAPI_Init(): Loaded local 'steamclient.so' OK.
CApplicationManagerPopulateThread took 0 milliseconds to initialize (will have waited on CAppInfoCacheReadFromDiskThread)
CAppInfoCacheReadFromDiskThread took 2 milliseconds to initialize
RecordSteamInterfaceCreation (PID 110): SteamGameServer013 /
RecordSteamInterfaceCreation (PID 110): SteamUtils010 /
Setting breakpad minidump AppID = 240
Looking up breakpad interfaces from steamclient
Calling BreakpadMiniDumpSystemInit
RecordSteamInterfaceCreation (PID 110): SteamGameServer013 /
RecordSteamInterfaceCreation (PID 110): SteamUtils010 /
RecordSteamInterfaceCreation (PID 110): SteamNetworking006 /
RecordSteamInterfaceCreation (PID 110): SteamGameServerStats001 /
RecordSteamInterfaceCreation (PID 110): STEAMHTTP_INTERFACE_VERSION003 /
RecordSteamInterfaceCreation (PID 110): STEAMINVENTORY_INTERFACE_V003 /
RecordSteamInterfaceCreation (PID 110): STEAMUGC_INTERFACE_VERSION014 /
RecordSteamInterfaceCreation (PID 110): STEAMAPPS_INTERFACE_VERSION008 /
Setting breakpad minidump AppID = 232330
No account token specified; logging into anonymous game server account.  (Use sv_setsteamaccount to login to a persistent account.)
RecordSteamInterfaceCreation (PID 110): SteamGameServer013 /
RecordSteamInterfaceCreation (PID 110): SteamUtils010 /
RecordSteamInterfaceCreation (PID 110): SteamNetworking006 /
RecordSteamInterfaceCreation (PID 110): SteamGameServerStats001 /
RecordSteamInterfaceCreation (PID 110): STEAMHTTP_INTERFACE_VERSION003 /
RecordSteamInterfaceCreation (PID 110): STEAMINVENTORY_INTERFACE_V003 /
RecordSteamInterfaceCreation (PID 110): STEAMUGC_INTERFACE_VERSION014 /
RecordSteamInterfaceCreation (PID 110): STEAMAPPS_INTERFACE_VERSION008 /
ConVarRef room_type doesn't point to an existing ConVar
Executing dedicated server config file server.cfg
Using map cycle file 'cfg/mapcycle_default.txt'.  ('cfg/mapcycle.txt' was not found.)
Set motd from file 'cfg/motd_default.txt'.  ('cfg/motd.txt' was not found.)
Set motd_text from file 'cfg/motd_text_default.txt'.  ('cfg/motd_text.txt' was not found.)
SV_ActivateServer: setting tickrate to 66.7
'server.cfg' not present; not executing.
Server is hibernating
Connection to Steam servers successful.
   Public IP is XX.YY.ZZ.AA.
Assigned anonymous gameserver Steam ID [X:Y:ZZZZZZZ:AAAA].
VAC secure mode is activated.
```

Check daemon shell:
```
# [shell DAEMON]

2023/08/02 19:18:12 Received request to update the game server
2023/08/02 19:18:12 Updating the game server...
2023/08/02 19:20:59 Received request to start the game server
2023/08/02 19:20:59 Starting game server...
2023/08/02 19:21:04 Received request to send game server logs
2023/08/02 19:21:04 Sending game server logs
2023/08/02 19:21:04 Sent 3712 bytes (60 lines) of game server logs
```

Issue a console command:
```
# [shell SERVER]

$ steamcmd-cli server console say hi
Command 'say hi' dispatched to game server console.

Console: hi
```

Check daemon shell:
```
# [shell DAEMON]

2023/08/02 19:33:11 Received server command 'say hi'
2023/08/02 19:33:11 Sending server command 'say hi' to game server console...
2023/08/02 19:33:12 Sending game server console replies for command 'say hi'...
```

View server logs again:
```
# [shell SERVER]

$ steamcmd-cli server logs
...
say hi
Console: hi
```

Stop the server:
```
# [shell SERVER]

$ steamcmd-cli server stop
Command 'quit' dispatched to game server console.
Game server stopped.
```

Check daemon shell:
```
# [shell DAEMON]

2023/08/02 19:35:18 Received request to stop the game server
2023/08/02 19:35:18 Received server command 'quit'
2023/08/02 19:35:18 Sending server command 'quit' to game server console...
```

Optionally the server could also be stopped by issuing the console command "quit":
```
# [shell SERVER]

$ steamcmd-cli server console quit
```

Now you can safely quit the daemon via CTRL+C. Quitting the daemon without stopping the server first will kill the server process and is therefore not encouraged. If you encounter any side effects after force-quitting, just update the server again and it should be fine. This would also be the case if you ran the server via `srcds_linux` or `srcds_run` directly.

# Installation

## Install from the released binaries
The tool is only available for the `i386` architecture. This is because Valve dedicated servers are written for it. Don't worry, it will run on `x86_64` systems as well.

Links: TODO (no releases yet)

## Install from source
**Note**: This step is currently not working at all.

If you have a working Go toolchain, you can use that to install `steamcmd-cli` to your `$GOPATH/bin`.

```
$ go install github.com/thetredev/steamcmd-cli@latest
```

## Build from source
Before building in any way, you'll need the sources:

```
$ git clone https://github.com/thetredev/steamcmd-cli.git
$ cd steamcmd-cli
```

### Locally
```
# Initialize the gopls workspace
$ go work init
$ go work use ./cli ./docker

# Initialize the dependencies
$ go get -C cli

# Build statically:
$ CGO_ENABLED=0 GOOS=linux go build -C cli -o $(pwd)/build/steamcmd-cli
```

### Build container images
This example uses Docker:
```
$ ./build_docker.sh
```

This command will create the following development images:
| Image | Description |
| --- | --- |
| `github.com/thetredev/steamcmd-cli:golang` | Official i386 `golang` v1.20 Alpine image with `git` installed. |
| `github.com/thetredev/steamcmd-cli:latest` | Scratch image with `steamcmd-cli` executable located at `/bin/steamcmd-cli`.<br/><br/>The entrypoint is set to `/bin/steamcmd-cli`. |
| `github.com/thetredev/steamcmd-cli:daemon` | Scratch image with `steamcmd-cli` executable located at `/bin/steamcmd-cli`. All necessary SteamCMD and HLDS/SRCDS dependencies are preinstalled. SteamCMD itself is preinstalled.<br/><br/>The entrypoint is set to `/bin/steamcmd daemon`. |

The containers can now be run locally:
```
# Daemon
$ docker run --rm -v $(pwd)/certs:/certs:ro github.com/thetredev/steamcmd-cli:daemon
2023/08/02 21:35:44 Listening for incoming requests on port 65000/TCP...

# CLI
$ docker run --rm -v $(pwd)/certs:/certs:ro github.com/thetredev/steamcmd-cli:latest
Custom SteamCMD client and game server manager implementation written in Go

Usage:
  steamcmd-cli [command]

Available Commands:
  certs       Certificate management for secure daemon/server communication
  completion  Generate the autocompletion script for the specified shell
  daemon      Run the daemon socket
  help        Help about any command
  server      Subcommands to communicate with the game server via the daemon socket

Flags:
  -h, --help      help for steamcmd-cli
  -v, --version   version for steamcmd-cli

Use "steamcmd-cli [command] --help" for more information about a command.
```

# Default configuration

See [./docs/environment-variables.md](./docs/environment-variables.md) for details.

# Certificate Configuration

To be able to communicate with the daemon, SSL certificate have to be setup first. You can use any valid CA and certificate/key pair combination, as long as the certificates were signed by the same CA.

This tool provides basic certificate management for getting everything up and running quickly. **BUT: You should always use CAs and certificate/key pairs that you trust!**.

## Generate certificates using this tool

First, create a CA:
```
$ steamcmd-cli certs ca ./certs/ca
```

Second, create the daemon certificate/key pair:
```
$ steamcmd-cli certs generate ./certs/ca/cert.pem ./certs/ca/cert.key ./certs/daemon
```

Third, create the server certificate/key pair:
```
$ steamcmd-cli certs generate ./certs/ca/cert.pem ./certs/ca/cert.key ./certs/server
```

Of course, you should provide [all the optional flags](./docs/commands.md#certs-command) to each command to be able to identify/verify the files.

## Daemon

To change the default values, either set the environment variables to different vales or pass the values to the executable, like this:

```
$ steamcmd-cli daemon /path/to/ca.pem /path/to/daemon.pem /path/to/daemon.key
```

For containers, you can simply put the paths in the correct directory structure and bind mount the path to `/certs`:

```
# Prepare the directory structure
$ mkdir -p /path/to/certs
$ mkdir -p /path/to/certs/ca
$ mkdir -p /path/to/certs/daemon

# Copy the files
$ cp ca.pem /path/to/certs/ca/cert.pem
$ cp daemon.pem /path/to/certs/daemon/cert.pem
$ cp daemon.key /path/to/certs/daemon/cert.key

# Run the container (using Docker in this example)
$ docker run -v /path/to/certs:/certs:ro ....
```

## Server

To change the default values, either set the environment variables to different vales or pass the values to the executable, like this:

```
$ steamcmd-cli server logs /path/to/server.pem /path/to/server.key
```

For containers, you can simply put the paths in the correct directory structure and bind mount the path to `/certs`:

```
# Prepare the directory structure
$ mkdir -p /path/to/certs/server

# Copy the files
$ cp daemon.pem /path/to/certs/server/cert.pem
$ cp daemon.key /path/to/certs/server/cert.key

# Run the container (using Docker in this example)
$ docker run -v /path/to/certs:/certs:ro ....
```

# Using container images
This section assumes container images were built using the `build_docker.sh` script.

## Create certificates

### CA
```
$ docker run --rm \
  -v $(pwd)/certs:/certs \
  --entrypoint /bin/steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest certs ca /certs/ca
```

### Daemon
```
$ docker run --rm \
  -v $(pwd)/certs:/certs \
  --entrypoint /bin/steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest certs generate /certs/ca/cert.pem /certs/ca/cert.key /certs/daemon
```

### Server
```
$ docker run --rm \
  -v $(pwd)/certs:/certs \
  --entrypoint /bin/steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest certs generate /certs/ca/cert.pem /certs/ca/cert.key /certs/server
```

### Fix local file permissions
```
$ sudo chown -R $(id -u):$(id -g) ./certs
```

## Create a container network for internal communication
```
$ docker network create steamcmd-cli
cd82a09ca68f10165ce52fbd8273fae4b3fd1201a9216336d2b91df74f71ebc0
```

## Run the Daemon container (Counter-Strike: Source)
```
$ docker run --rm \
  -v $(pwd)/certs:/certs:ro \
  -v /tmp/srcds/css:/steamcmd/server \
  -p 27015:27015/udp \
  -e STEAMCMD_SERVER_MAP=de_dust2 \
  -e STEAMCMD_SERVER_MAXPLAYERS=16 \
  -e STEAMCMD_SERVER_APPID=232330 \
  --network steamcmd-cli \
  --name daemon \
  github.com/thetredev/steamcmd-cli:daemon
```

## Manage the game server using the CLI container

### Update
```
$ docker run --rm \
  -v $(pwd)/certs:/certs:ro \
  -e STEAMCMD_CLI_SOCKET_IP=daemon \
  --network steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest server update
```

### Start
```
$ docker run --rm \
  -v $(pwd)/certs:/certs:ro \
  -e STEAMCMD_CLI_SOCKET_IP=daemon \
  --network steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest server start
```

### Logs
```
$ docker run --rm \
  -v $(pwd)/certs:/certs:ro \
  -e STEAMCMD_CLI_SOCKET_IP=daemon \
  --network steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest server logs
```

### Console
```
$ docker run --rm \
  -v $(pwd)/certs:/certs:ro \
  -e STEAMCMD_CLI_SOCKET_IP=daemon \
  --network steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest server console say hi
```

### Stop
```
$ docker run --rm \
  -v $(pwd)/certs:/certs:ro \
  -e STEAMCMD_CLI_SOCKET_IP=daemon \
  --network steamcmd-cli \
  github.com/thetredev/steamcmd-cli:latest server stop
```

## Or: Use the CLI tool locally to interact with the Daemon container

### Run the Daemon container with an exposed port (Counter-Strike: Source)
```
$ docker run --rm \
  -v $(pwd)/certs:/certs:ro \
  -v /tmp/srcds/css:/steamcmd/server \
  -p 27015:27015/udp \
  -p 65000:65000/tcp \
  -e STEAMCMD_SERVER_MAP=de_dust2 \
  -e STEAMCMD_SERVER_MAXPLAYERS=16 \
  -e STEAMCMD_SERVER_APPID=232330 \
  github.com/thetredev/steamcmd-cli:daemon
```

### Run the CLI locally
```
$ export STEAMCMD_CLI_SOCKET_IP=127.0.0.1
$ export STEAMCMD_CLI_CERTIFICATE_PATH=$(pwd)/certs/server/cert.pem
$ export STEAMCMD_CLI_CERTIFICATE_KEY_PATH=$(pwd)/certs/server/cert.key
$ ./steamcmd-cli server update
...
```

# Shell auto-completions

### Bash
System-wide:
```
$ steamcmd-cli completion > /etc/bash_completion.d/steamcmd-cli
```

User-specific:
```
$ echo 'source <(steamcmd-cli completion bash)' >> ~/.bashrc
```

### Zsh
System-wide:
```
$ steamcmd-cli completion > /usr/local/share/zsh/site-functions/_steamcmd-cli
```

User-specific:
```
$ echo 'source <(steamcmd-cli completion zsh)' >> ~/.zshrc
```

### Fish
User-specific:
```
$ steamcmd-cli completion > ~/.config/fish/completions/steamcmd-cli.fish
```
