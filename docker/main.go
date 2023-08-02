package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

func getUserID(idString string) int {
	id, err := strconv.Atoi(idString)

	if err != nil {
		log.Fatal(err)
	}

	return id
}

func getUserIDs(user *user.User) (int, int) {
	return getUserID(user.Uid), getUserID(user.Gid)
}

func main() {
	user, err := user.Lookup("steamcmd")

	if err != nil {
		log.Fatal(err)
	}

	uid, gid := getUserIDs(user)
	serverHome := os.Getenv("STEAMCMD_SERVER_HOME")

	err = os.MkdirAll(serverHome, os.FileMode(0770))

	if err != nil {
		fmt.Printf("mkdir: %s\n", serverHome)
		log.Fatal(err)
	}

	err = filepath.Walk(serverHome, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = syscall.Chown(name, uid, gid)
		}

		return err
	})

	if err != nil {
		fmt.Println("walk")
		log.Fatal(err)
	}

	err = os.MkdirAll("/tmp", os.ModePerm)

	if err != nil {
		fmt.Println("mkdir: /tmp")
		log.Fatal(err)
	}

	err = os.Chown("/tmp", uid, gid)

	if err != nil {
		fmt.Println("chown: tmp")
		log.Fatal(err)
	}

	command := exec.Command("/bin/steamcmd-cli", "daemon")
	command.Env = os.Environ()
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.SysProcAttr = &syscall.SysProcAttr{
		Pgid:    1,
		Setpgid: true,
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	err = command.Run()

	if err != nil {
		log.Fatal(err)
	}
}
