/**
 * @Time : 2022/12/29 9:49 AM
 * @Author : solacowa@gmail.com
 * @File : service_camera
 * @Software: GoLand
 */

package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var (
	cameraCmd = &cobra.Command{
		Use:               "camera command <args> [flags]",
		Short:             "相机操作命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## 相机操作命令
可用的配置类型：
[jpeg]

pi-kit camera -h
`,
	}

	cameraJpegCmd = &cobra.Command{
		Use:               "jpeg <args> [flags]",
		Short:             "拍照并生成jpeg文件",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## 拍照并生成jpeg文件
可用的配置类型：

pi-kit camera jpeg -h
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			return cameraJpegOutput(cmd.Context(), cameraJpegOutputPath)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepare(cmd.Context())
		},
	}
)

func cameraJpegOutput(ctx context.Context, outputPath string) (err error) {
	logger = log.With(logger, "cameraJpegCmd", "cameraJpegOutput", "outputPath", outputPath)
	if strings.EqualFold(PiConnect(piConnect).String(), PiConnectRemote.String()) {
		host := fmt.Sprintf("%s:%d", piHost, piSSHPort)

		//pKey, _ := ioutil.ReadFile("~/.ssh/id_rsa")
		//pKey := []byte("<privateKey>")
		//fmt.Println(string(pKey))

		//var signer ssh.Signer
		//signer, err = ssh.ParsePrivateKey(pKey)
		//if err != nil {
		//	_ = level.Warn(logger).Log("ssh", "ParsePrivateKey", "err", err.Error())
		//}

		var hostkeyCallback ssh.HostKeyCallback
		hostkeyCallback, err = knownhosts.New("~/.ssh/known_hosts")
		if err != nil {
			_ = level.Warn(logger).Log("knownhosts", "New", "err", err.Error())
		}

		conf := &ssh.ClientConfig{
			User:            piUser,
			HostKeyCallback: hostkeyCallback,
			Auth: []ssh.AuthMethod{
				ssh.Password(piPassword),
				//ssh.PublicKeys(signer),
			},
		}
		var conn *ssh.Client

		conn, err = ssh.Dial("tcp", host, conf)
		if err != nil {
			_ = level.Error(logger).Log("ssh", "Dial", "err", err.Error())
			return err
		}
		defer func(conn *ssh.Client) {
			_ = conn.Close()
		}(conn)

		var session *ssh.Session
		var stdin io.WriteCloser
		var stdout, stderr io.Reader

		session, err = conn.NewSession()
		if err != nil {
			_ = level.Error(logger).Log("conn", "NewSession", "err", err.Error())
			return err
		}
		defer func(session *ssh.Session) {
			_ = session.Close()
		}(session)

		stdin, err = session.StdinPipe()
		if err != nil {
			_ = level.Error(logger).Log("session", "StdinPipe", "err", err.Error())
			return err
		}

		stdout, err = session.StdoutPipe()
		if err != nil {
			_ = level.Error(logger).Log("session", "StdoutPipe", "err", err.Error())
			return err
		}

		stderr, err = session.StderrPipe()
		if err != nil {
			_ = level.Error(logger).Log("session", "StderrPipe", "err", err.Error())
			return err
		}

		wr := make(chan []byte, 10)

		go func() {
			for {
				select {
				case d := <-wr:
					_, err := stdin.Write(d)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stdout)
			for {
				if tkn := scanner.Scan(); tkn {
					rcv := scanner.Bytes()

					raw := make([]byte, len(rcv))
					copy(raw, rcv)

					fmt.Println(string(raw))
				} else {
					if scanner.Err() != nil {
						fmt.Println(scanner.Err())
					} else {
						fmt.Println("io.EOF")
					}
					return
				}
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)

			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()

		if err = session.Shell(); err != nil {
			_ = level.Error(logger).Log("session", "Shell", "err", err.Error())
		}

		for {
			fmt.Println("$")

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			text := scanner.Text()

			wr <- []byte(text + "\n")
		}
	}
	return
}
