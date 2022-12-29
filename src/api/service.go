/**
 * @Time : 2022/12/29 9:37 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package api

import (
	"context"
	"github.com/icowan/pi-kit/src/api/ssh"
)

type Service interface {
	SSHClient(ctx context.Context) ssh.Service
}

type service struct {
	sshClient ssh.Service
}

func (s *service) SSHClient(ctx context.Context) ssh.Service {
	return s.sshClient
}

func New() Service {
	sshClient := ssh.New()
	return &service{sshClient: sshClient}
}
