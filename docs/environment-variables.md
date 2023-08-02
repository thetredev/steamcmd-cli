# Overview of environment variables and their default values

The default values were chosen because the primary way of running the daemon and server is using a container runtime such as Docker.

## Daemon

These only apply to the `daemon` command.

| Environment variable | Description | Default value |
| --- | --- | --- |
| `STEAMCMD_SH` | Path to the `steamcmd` executable. | Container image: `/steamcmd/package/steamcmd.sh`<br/>Empty otherwise. |
| `STEAMCMD_SERVER_APPID` | The dedicated server app ID. | Empty |
| `STEAMCMD_SERVER_MOD` | The dedicated server mod. Used to apply mods to app ID `90`. Example: `czero` for Counter-Strike: Condition Zero. | Empty  |
| `STEAMCMD_SERVER_HOME` | The path to where the dedicated server files are downloaded and the dedicated server is run from. | Container image: `/steamcmd/server`<br/>Empty otherwise. |
| `STEAMCMD_SERVER_PORT` | The dedicated server port. Note: if running inside a container, this **must** be the same for both the container and host. At least for HLDS/SRCDS servers.<br/><br/>`docker run -p <port>:<port> ....` | Container image: `27015`<br/>Empty otherwise. |
| `STEAMCMD_SERVER_MAXPLAYERS` | Amount of available player slots. | Empty |
| `STEAMCMD_SERVER_MAP` | The map to start the server with. | Empty |
| `STEAMCMD_SERVER_TICKRATE` | The tickrate to start the server with. Only applies to Counter-Strike: Global Offensive and other dedicated servers based on the CS:GO engine (as far as I'm aware). If you want other servers to have a higher tickrate than their official maximum, you must take care of tickrate enablers yourself. | `128` |
| `STEAMCMD_SERVER_THREADS` | The amount of threads the server uses. | `3`
| `STEAMCMD_SERVER_FPSMAX` | The `fps_max` value to start the server with. | `300`
| `STEAMCMD_DAEMON_CA_CERTIFICATE_PATH` | Daemon CA certificate path | `/certs/ca/cert.pem` |
| `STEAMCMD_DAEMON_CERTIFICATE_PATH` | Daemon certificate path | `/certs/daemon/cert.pem` |
| `STEAMCMD_DAEMON_CERTIFICATE_KEY_PATH` | Daemon certificate key path | `/certs/daemon/cert.key`

## Server

These only apply to the `server` command.

| Environment variable | Description | Default value |
| --- | --- | --- |
| `STEAMCMD_CLI_CERTIFICATE_PATH` | Server certificate path | `/certs/server/cert.pem` |
| `STEAMCMD_CLI_CERTIFICATE_KEY_PATH` | Server certificate key path | `/certs/server/cert.key`

## Shared

These can be applied to both the `daemon` and the `server` command.

| Environment variable | Description | Default value |
| --- | --- | --- |
| `STEAMCMD_CLI_SOCKET_IP` | The IP address or hostname of the TCP socket.<br/><br/>For the daemon, this value is used to identify the network interface to listen on. `0.0.0.0` means *listen on all interfaces* in this case.<br/><br/>For the server, this is the destination address of the daemon. Cannot be `0.0.0.0`. If the daemon is running on the same machine, use `127.0.0.1` or `localhost`. | `0.0.0.0` |
| `STEAMCMD_CLI_SOCKET_PORT` | The port of the TCP socket.<br/><br/>For the daemon, this value is used as the port to open the socket on.<br/><br/>For the server, this is the destination port of the daemon.<br/><br/>In either case, if running as an unprivileged user, this value cannot be below `1024`. See https://www.w3.org/Daemon/User/Installation/PrivilegedPorts.html for an explanation. It also cannot be higher than `65535` due to a limiation of the TCP specification (16 bits). | `65000` |
