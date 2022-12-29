/**
 * @Time: 2020/2/14 13:57
 * @Author: solacowa@gmail.com
 * @File: encode
 * @Software: GoLand
 */

package encode

import (
	"context"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
)

type Response struct {
	Success     bool        `json:"success"`
	Code        int         `json:"code"`
	Data        interface{} `json:"data,omitempty"`
	Error       error       `json:"-"`
	Message     string      `json:"message,omitempty"`
	TraceId     string      `json:"traceId,omitempty"`
	HttpHeaders http.Header `json:"-"`
	StatusCode  int         `json:"-"`
	Body        []byte      `json:"-"`
}

type Failure interface {
	Failed() error
}

type Errorer interface {
	Error() error
}

func Error(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func JsonError(ctx context.Context, err error, w http.ResponseWriter) {
	var errDefined bool
	for k := range ResponseMessage {
		if strings.Contains(err.Error(), k.Error().Error()) {
			errDefined = true
			break
		}
	}

	if !errDefined {
		err = errors.Wrap(ErrSystem.Error(), err.Error())
	}
	if err == nil {
		err = errors.Wrap(err, ErrSystem.Error().Error())
	}
	traceId, _ := ctx.Value("traceId").(string)
	w.Header().Set("TraceId", traceId)
	_ = kithttp.EncodeJSONResponse(ctx, w, map[string]interface{}{
		"message": err.Error(),
		"code":    ResponseMessage[ResStatus(strings.Split(err.Error(), ":")[0])],
		"success": false,
		"traceId": traceId,
	})
}

func JsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(Failure); ok && f.Failed() != nil {
		JsonError(ctx, f.Failed(), w)
		return nil
	}
	resp := response.(Response)
	if resp.Error == nil {
		resp.Code = 200
		resp.Success = true
	} else {
		resp.Code = ResponseMessage[ResStatus(strings.Split(resp.Error.Error(), ":")[0])]
		resp.Message = resp.Error.Error()
	}
	traceId, _ := ctx.Value("traceId").(string)
	resp.TraceId = traceId
	w.Header().Set("TraceId", traceId)
	return kithttp.EncodeJSONResponse(ctx, w, resp)
}

func (r Response) Headers() http.Header {
	return r.HttpHeaders
}

func (r Response) StatusCoder() int {
	return r.StatusCode
}
