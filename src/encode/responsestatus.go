/**
 * @Time: 2020/3/27 17:34
 * @Author: solacowa@gmail.com
 * @File: responsestatus
 * @Software: GoLand
 */

package encode

import (
	"github.com/pkg/errors"
)

type ResStatus string

var ResponseMessage = map[ResStatus]int{
	Invalid:            400,
	InvalidParams:      400,
	InvalidParamsAlias: 400,
	InvalidParamsName:  400,
	ErrParamsPhone:     401,
	ErrBadRoute:        401,
	ErrSystem:          500,
	ErrNotfound:        404,
	ErrLimiter:         429,
}

const (
	// 公共错误信息
	Invalid            ResStatus = "invalid"
	InvalidParams      ResStatus = "请求参数错误"
	ErrNotfound        ResStatus = "不存在"
	ErrBadRoute        ResStatus = "请求路由错误"
	ErrParamsPhone     ResStatus = "手机格式不正确"
	ErrLimiter         ResStatus = "太快了,等我一会儿..."
	InvalidParamsAlias ResStatus = "别名格式不正确"
	InvalidParamsName  ResStatus = "名称不合规: ^[a-z]([-a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"
	ErrSystem          ResStatus = "系统错误"
)

func (c ResStatus) String() string {
	return string(c)
}

func (c ResStatus) Error() error {
	return errors.New(string(c))
}

func (c ResStatus) Wrap(err error) error {
	return errors.Wrap(err, string(c))
}
