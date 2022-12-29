/**
 * @Time: 2022/12/29 22:45
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package ssh

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"net"
)

type Service interface {
	Connection(ctx context.Context, host, user, password string) (res *ssh.Client, err error)
}

type service struct {
}

func (s *service) Connection(ctx context.Context, host, user, password string) (res *ssh.Client, err error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "ssh.Dial")
	}
	return conn, nil
}

func New() Service {
	return &service{}
}
