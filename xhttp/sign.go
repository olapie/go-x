package xhttp

import (
	"net/http"

	"go.olapie.com/x/xbase62"
	"go.olapie.com/x/xcontext"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xlog"
)

func SignRequest(req *http.Request, createAPIKey func(h http.Header)) {
	a := xcontext.GetOutgoingActivity(req.Context())
	if a != nil {
		xcontext.CopyActivityHeader(req.Header, a)
	} else {
		xlog.FromContext(req.Context()).Warn("no outgoing context")
	}
	if traceID := xhttpheader.GetTraceID(req.Header); traceID == "" {
		if traceID == "" {
			traceID = xbase62.NewUUIDString()
			xlog.FromContext(req.Context()).Info("xhttp.SignRequest: generated trace id " + traceID)
		}
		xhttpheader.SetTraceID(req.Header, traceID)
	}

	if createAPIKey == nil {
		createAPIKey = xhttpheader.Sign[http.Header]
	}
	createAPIKey(req.Header)
}
