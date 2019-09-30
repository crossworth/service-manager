package servicemanager

import (
	"net/http/httptest"
	"testing"
)

func TestGetRealIP(t *testing.T) {
	want := "138.1.2.3"

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Real-Ip", want)

	got := getRealIP(req)
	if got != want {
		t.Fatalf("got wrong, got %q want %q", got, want)
	}
}

func TestGetMD5Hash(t *testing.T) {
	want := "6980afb078fd03156924b050c1df180a"
	got := md5Hash("my-md5-hash")
	if got != want {
		t.Fatal("could not hash md5")
	}
}
