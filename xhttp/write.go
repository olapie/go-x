package xhttp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"go.olapie.com/x/xerror"
	"go.olapie.com/x/xhttpheader"
	"go.olapie.com/x/xmime"
)

func BasicAuthenticate(w http.ResponseWriter, realm string) {
	a := "Basic realm=" + strconv.Quote(realm)
	w.Header().Set(xhttpheader.KeyWWWAuthenticate, a)
	w.WriteHeader(http.StatusUnauthorized)
}

func Error(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if errors.Is(err, xerror.DBNoRecords) || errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	status := http.StatusInternalServerError
	if s := xerror.GetCode(err); s != 0 {
		status = s
	}

	if status < 100 || status > 599 {
		log.Println("invalid status:", status)
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	_, _ = w.Write([]byte(err.Error()))
}

func JSON(w http.ResponseWriter, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	xhttpheader.SetContentType(w.Header(), xmime.JsonUTF8)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

func JSONOrError(w http.ResponseWriter, v any, err error) {
	if err != nil {
		Error(w, err)
	} else {
		JSON(w, v)
	}
}

func OctetStream(w http.ResponseWriter, b []byte) {
	xhttpheader.SetContentType(w.Header(), xmime.OctetStream)
	_, err := w.Write(b)
	if err != nil {
		slog.Error("cannot write", "err", err.Error())
	}
}

func HTMLFile(w http.ResponseWriter, s string) {
	xhttpheader.SetContentType(w.Header(), xmime.HtmlUTF8)
	_, err := w.Write([]byte(s))
	if err != nil {
		slog.Error("cannot write", "err", err.Error())
	}
}

func CSSFile(w http.ResponseWriter, s string) {
	xhttpheader.SetContentType(w.Header(), xmime.CSS)
	_, err := w.Write([]byte(s))
	if err != nil {
		slog.Error("cannot write", "err", err.Error())
	}
}

func JSFile(w http.ResponseWriter, s string) {
	xhttpheader.SetContentType(w.Header(), xmime.Javascript)
	_, err := w.Write([]byte(s))
	if err != nil {
		slog.Error("cannot write", "err", err.Error())
	}
}

func StreamFile(w http.ResponseWriter, name string, f io.ReadCloser) {
	defer f.Close()
	xhttpheader.SetContentType(w.Header(), xmime.OctetStream)
	if name != "" {
		w.Header().Set(xhttpheader.KeyContentDisposition, xhttpheader.ToAttachment(name))
	}
	_, err := io.Copy(w, f)
	if err != nil {
		if err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
