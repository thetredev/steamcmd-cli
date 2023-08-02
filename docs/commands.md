# CLI commands overview

## main command
```
$ steamcmd-cli -h
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

For the full list of subcommands, see [Usage](#usage).

## certs command
```
$ steamcmd-cli certs -h
Certificate management for secure daemon/server communication

Usage:
  steamcmd-cli certs [flags]
  steamcmd-cli certs [command]

Available Commands:
  ca          Generate CA certificate and private key
  generate    Generate daemon/server certificate and private key

Flags:
  -h, --help      help for certs
  -v, --version   version for certs

Use "steamcmd-cli certs [command] --help" for more information about a command.
```

### certs ca command
```
$ steamcmd-cli certs ca -h
Generate CA certificate and private key

Usage:
  steamcmd-cli certs ca <output dir> [flags]

Flags:
      --company string          certificate company (default "steamcmd-cli")
      --country string          certificate country (default "US")
      --hostname stringArray    certificate hostname list.
                                DNS names and IPv4/IPv6 addresses are valid, IP networks are NOT.
      --key-length int          certificate key length (default 4096)
      --locality string         certificate locality (default "server1")
      --postal-code string      certificate postal code (default "row1")
      --province string         certificate province (default "datacenter")
      --street-address string   certificate street address (default "rack1")
      --valid-until string      Date until the certificate should be valid (UTC).
                                Hours, minutes, and seconds will always be '01'. (default "2999-12-31")
  -h, --help                    help for ca
  -v, --version                 version for ca
```

### certs generate command
```
$ steamcmd-cli certs generate -h
Generate daemon/server certificate and private key

Usage:
  steamcmd-cli certs generate <CA certificate input file path> <CA certificate input key path> <output dir> [flags]

Flags:
      --company string          certificate company (default "steamcmd-cli")
      --country string          certificate country (default "US")
      --hostname stringArray    certificate hostname list.
                                DNS names and IPv4/IPv6 addresses are valid, IP networks are NOT.
      --key-length int          certificate key length (default 4096)
      --locality string         certificate locality (default "server1")
      --postal-code string      certificate postal code (default "row1")
      --province string         certificate province (default "datacenter")
      --street-address string   certificate street address (default "rack1")
      --valid-until string      Date until the certificate should be valid (UTC).
                                Hours, minutes, and seconds will always be '01'. (default "2999-12-31")
  -h, --help                    help for generate
  -v, --version                 version for generate
```

## daemon command
```
$ steamcmd-cli daemon -h
Run the daemon socket

Usage:
  steamcmd-cli daemon <CA certificate input file path> <Daemon certificate input key path> [flags]

Flags:
  -h, --help      help for daemon
  -v, --version   version for daemon
```

## server command
```
$ steamcmd-cli server -h
Subcommands to communicate with the game server via the daemon socket

Usage:
  steamcmd-cli server [flags]
  steamcmd-cli server [command]

Available Commands:
  console     Send commands to the game server console via the daemon socket
  logs        Retrieve game server logs from the daemon socket
  start       Start the game server via the daemon socket
  stop        Stop the game server via the daemon socket
  update      Update the game server via the daemon socket

Flags:
  -h, --help      help for server
  -v, --version   version for server

Use "steamcmd-cli server [command] --help" for more information about a command.
```

### server console command
```
$ steamcmd-cli server console -h
Send commands to the game server console via the daemon socket

Usage:
  steamcmd-cli server console [server certificate file] [server key file] <command list> [flags]

Flags:
  -h, --help      help for console
  -v, --version   version for console
```

### server logs command
```
$ steamcmd-cli server logs
Retrieve game server logs from the daemon socket

Usage:
  steamcmd-cli server logs [server certificate file] [server key file] [flags]

Flags:
  -h, --help      help for logs
  -v, --version   version for logs
```

### server start command
```
$ steamcmd-cli start -h
Start the game server via the daemon socket

Usage:
  steamcmd-cli server start [server certificate file] [server key file] [flags]

Flags:
  -h, --help      help for start
  -v, --version   version for start
```

### server stop command
```
$ steamcmd-cli stop -h
Stop the game server via the daemon socket

Usage:
  steamcmd-cli server stop [server certificate file] [server key file] [flags]

Flags:
  -h, --help      help for stop
  -v, --version   version for stop
```

### server update command
```
$ steamcmd-cli update -h
Update the game server via the daemon socket

Usage:
  steamcmd-cli server update [server certificate file] [server key file] [flags]

Flags:
  -h, --help      help for update
  -v, --version   version for update
```
