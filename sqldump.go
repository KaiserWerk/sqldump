package sqldump

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Uploader contains credentials for a Hetzner storage box. Make sure to enable external access.
type Uploader struct {
	Username string
	Password string
	Host     string

	counter  uint8
	maxFiles uint8
}

// NewUploader creates a new Uploader instance with the given credentials.
func NewUploader(username, password, host string) *Uploader {
	return &Uploader{
		Username: username,
		Password: password,
		Host:     host,
	}
}

// ScheduleUpload schedules periodic uploads of the specified file to the storage box at the given interval.
func (u *Uploader) ScheduleUpload(ctx context.Context, filenameFunc func() (string, error), deleteFileAfterUpload bool, targetDirectory string, interval time.Duration, maxFiles uint8) {
	u.maxFiles = maxFiles
	u.counter = 0
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				ticker = nil
				return
			case <-ticker.C:
				filename, err := filenameFunc()
				if err != nil {
					log.Printf("Error getting filename: %v", err)
					continue
				}
				if err := u.UploadFileToStorageBox(filename, u.getNextFilename(filename), targetDirectory); err != nil {
					log.Printf("Error uploading file: %v", err)
				}

				_ = os.Remove(filename)
			}
		}
	}()
}

func (u *Uploader) getNextFilename(filename string) string {
	name := filepath.Base(filename)
	nameOnly := name[:len(name)-len(filepath.Ext(name))]
	ext := filepath.Ext(filename)
	next := fmt.Sprintf("%s-%d%s", nameOnly, u.incrementCounter(), ext)
	return next
}

func (u *Uploader) incrementCounter() uint8 {
	if u.counter < u.maxFiles {
		u.counter++
		return u.counter
	}

	// reset
	u.counter = 1
	return u.counter
}

func (u *Uploader) UploadFileToStorageBox(filename, targetFilename string, targetDirectory string) error {
	// SSH Config
	config := &ssh.ClientConfig{
		User: u.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(u.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // haha
	}

	// ssh connection
	conn, err := ssh.Dial("tcp", u.Host, config)
	if err != nil {
		return fmt.Errorf("SSH Dial error: %w", err)
	}
	defer conn.Close()

	// create new SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("SFTP error: %w", err)
	}
	defer client.Close()

	srcFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}
	defer srcFile.Close()

	// target on storage box
	dstFile, err := client.Create(targetDirectory + "/" + targetFilename)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	defer dstFile.Close()

	// copy file content
	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return fmt.Errorf("upload error: %w", err)
	}

	return nil
}
