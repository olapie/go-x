package xsecurity

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"go.olapie.com/x/xnet"
)

func TestNewTLSConfig(t *testing.T) {
	config, err := NewTLSConfig()
	if err != nil {
		t.Fatal(err)
	}

	port := xnet.FindTCPPort("127.0.0.1", 8000, 9000)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	server := http.Server{
		Addr:                         addr,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    config,
	}
	server.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})
	go func() {
		err = server.ListenAndServeTLS("", "")
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				t.Error(err)
			}
			return
		}
	}()

	time.Sleep(1 * time.Second)
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get("https://" + addr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "ok" {
		t.Fatal("body is " + string(body))
	}
	_ = server.Shutdown(context.Background())
}
