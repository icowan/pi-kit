/**
 * @Time: 2022/12/29 22:45
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package ssh

import (
	"context"
)

type Service interface {
	Connection(ctx context.Context, host, user, password, keyFile string) (res *Client, err error)
}

type service struct {
}

func (s *service) Connection(ctx context.Context, host, user, password, keyFile string) (res *Client, err error) {
	client, err := DialWithKey(host, user, keyFile)
	if err != nil {
		return DialWithPasswd(host, user, password)
	}
	return client, nil
}

func New() Service {
	return &service{}
}
