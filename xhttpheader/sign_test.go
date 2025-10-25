package xhttpheader

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestCreateSimpleAPIKey(t *testing.T) {
	header := http.Header{}
	SetClientID(header, uuid.NewString())
	SetTraceID(header, uuid.NewString())
	Sign(header)
	fmt.Println(header)
	time.Sleep(2 * time.Second)
	if err := Verify(header, time.Second*5); err != nil {
		t.Fatal(err)
	}

	if err := Verify(header, time.Second); err == nil {
		t.Fatal("should have expired")
	}
}
