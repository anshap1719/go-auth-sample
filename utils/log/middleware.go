package log

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/goadesign/goa"
)

func LogContextMiddleware(h goa.Handler) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		if l, ok := goa.ContextLogger(ctx).(*adapter); ok {
			newLogger := l.SetContext(ctx)
			ctx = goa.WithLogger(ctx, newLogger)
		}
		return h(ctx, rw, req)
	}
}

func LogInternalError(h goa.Handler) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		resp := httptest.NewRecorder()
		respData := goa.ContextResponse(ctx)
		old := respData.SwitchWriter(resp)
		err := h(ctx, resp, req)
		if err != nil {
			return err
		}

		respData.SwitchWriter(old)

		for k, sl := range resp.HeaderMap {
			for _, v := range sl {
				rw.Header().Add(k, v)
			}
		}

		if resp.Code != 0 {
			rw.WriteHeader(resp.Code)
		}

		if resp.Code == 500 {
			Error(ctx, "Internal Server Error: %s", resp.Body.String())
			err = json.NewEncoder(rw).Encode(goa.ErrInternal("An error has occurred, please reload and try again"))
			if err != nil {
				return err
			}
		} else {
			_, err = io.Copy(rw, resp.Body)
		}

		return err
	}
}
