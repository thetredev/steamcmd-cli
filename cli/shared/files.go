package shared

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func calculateFileBufferSize(fileSize int64) int64 {
	const maxBufferSize = 1024 * 1024 * 8
	bufferSize := fileSize

	for {
		if bufferSize <= maxBufferSize {
			break
		}

		bufferSize /= 2
	}

	return bufferSize
}

func TransferFile(socket SharedSocket, rootPath string, filePath string, tcpDelay time.Duration, logger *log.Logger) error {
	file, err := os.Open(filepath.Join(rootPath, filePath))

	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	socket.SendMessage(fmt.Sprintf("%s;%s;%d", MESSAGE_SERVER_FILE_TRANSFER_START, filePath, fileSize))
	time.Sleep(tcpDelay)

	bufferSize := calculateFileBufferSize(fileSize)

	buffer := make([]byte, bufferSize)
	bytesSent := int64(0)

	for {
		n, err := file.ReadAt(buffer, bytesSent)

		if n == 0 {
			if err != nil && err != io.EOF {
				return err
			}

			break
		}

		socket.SendMessage(fmt.Sprintf("%s;%d", MESSAGE_SERVER_FILE_TRANSFER_LEN, n))
		time.Sleep(tcpDelay)

		socket.SendBytes(buffer)
		time.Sleep(tcpDelay)

		bytesSent += int64(n)
		logger.Printf("Transferred %d (%d/%d) bytes over the socket.", n, bytesSent, fileSize)
	}

	socket.SendMessage(MESSAGE_SERVER_FILE_TRANSFER_END)
	time.Sleep(tcpDelay)

	return file.Close()
}

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

func ReceiveFile(reader *bufio.Reader, filePath string, fileSize int64) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0775)

	if err != nil {
		return fmt.Errorf("ERR MKDIR: %s", err.Error())
	}

	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("ERR CREATE FILE: %s", err.Error())
	}

	for {
		count, buffer, err := ReadBuffer(reader, 256)

		if err != nil {
			fmt.Printf("ERR SOCKET READ DL: %s\n", err.Error())
			break
		}

		if count == 0 {
			break
		}

		message := strings.Split(string(buffer[:count]), "\n")[0]

		if strings.HasPrefix(message, MESSAGE_SERVER_FILE_TRANSFER_END) {
			fmt.Println("")
			break
		}

		if strings.Contains(message, MESSAGE_SERVER_FILE_TRANSFER_LEN) {
			message = strings.Split(message, ";")[1]
			bufferSizeParsed, err := strconv.ParseInt(message, 10, 64)

			if err != nil {
				return fmt.Errorf("ERR BUFFER SIZE CONV: %s", err.Error())
			}

			bufferSize := int(bufferSizeParsed)

			progressBar := progressbar.NewOptions(bufferSize,
				progressbar.OptionFullWidth(),
				progressbar.OptionShowBytes(true),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}),
			)

			for !progressBar.IsFinished() {
				count, buffer, err := ReadBuffer(reader, bufferSize)

				if err != nil {
					fmt.Printf("DL INNER ERR: %s\n", err.Error())
					break
				}

				if count == 0 {
					break
				}

				n, err := file.Write(buffer[:count])

				if err != nil {
					fmt.Printf("ERROR WRITE FILE: %s\n", err.Error())
					break
				} else {
					progressBar.Add(n)
				}
			}
		}
	}

	return file.Close()
}
