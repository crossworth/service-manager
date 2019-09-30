package client

import (
	"testing"
)

func TestRegister(t *testing.T) {

}

func TestGetLocalIP(t *testing.T) {
	want := "192.168.0.10"
	got, err := GetLocalIP()
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf("got wrong IP, want %q got %q", want, got)
	}
}