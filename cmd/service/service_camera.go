/**
 * @Time : 2022/12/29 9:49 AM
 * @Author : solacowa@gmail.com
 * @File : service_camera
 * @Software: GoLand
 */

package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"
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
			return cameraJpegOutputExec(cmd.Context(), cameraJpegOutput)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prepare(cmd.Context())
		},
	}
)

func cameraJpegOutputExec(ctx context.Context, outputPath string) (err error) {
	if strings.EqualFold(outputPath, "") {
		outputPath, _ = os.Getwd()
	}
	outputPath = fmt.Sprintf("%s/camera-%s.jpeg", outputPath, time.Now().Format("20060102150405"))
	logger = log.With(logger, "cameraJpegCmd", "cameraJpegOutput", "outputPath", outputPath)
	if strings.EqualFold(PiConnect(piConnect).String(), PiConnectRemote.String()) {
		cmds := []string{
			cmdCameraJpegBin,
			"--output",
			cameraJpegOutput,
		}
		cmd := strings.Join(cmds, " ")
		_ = level.Info(logger).Log("cmd", cmd)

		host := fmt.Sprintf("%s:%d", piHost, piSSHPort)
		sshCli, err := apiSvc.SSHClient(ctx).Connection(ctx, host, piUser, piPassword, "~/.ssh/id_rsa")
		if err != nil {
			_ = level.Error(logger).Log("apiSvc.SSHClient", "Connection", "err", err.Error())
			return errors.Wrap(err, "apiSvc.SSHClient.Connection")
		}
		defer func() {
			if closeErr := sshCli.Close(); closeErr != nil {
				_ = level.Warn(logger).Log("sshCli", "Close", "err", err.Error())
			}
		}()
		outBytes, err := sshCli.Cmd(cmd).Output()
		if err != nil {
			_ = level.Error(logger).Log("sshCli.Cmd", "Output", "err", err.Error())
			return errors.Wrap(err, "sshCli.Cmd.Output")
		}

		fmt.Println(string(outBytes))
		return nil
	}
	outBytes, err := exec.CommandContext(ctx, cmdCameraJpegBin, "-output", outputPath).Output()
	if err != nil {
		_ = level.Error(logger).Log("exec", "CommandContext", "err", err.Error())
		return errors.Wrap(err, "exec.CommandContext")
	}
	_ = level.Info(logger).Log("exec.CommandContext", "success")
	fmt.Println(string(outBytes))
	return
}
