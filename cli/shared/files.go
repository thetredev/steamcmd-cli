package shared

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func TransferFileList(socket SharedSocket, rootPath string, walkPath string, tcpDelay time.Duration) (int, error) {
	filesFound := 0

	err := filepath.Walk(walkPath, func(name string, info os.FileInfo, err error) error {
		if name == walkPath && info != nil && info.IsDir() {
			return nil
		}

		if err == nil {
			socket.SendMessage(strings.Split(name, rootPath)[1][1:])
			time.Sleep(tcpDelay)

			filesFound++
		}

		return err
	})

	return filesFound, err
}
