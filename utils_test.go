package servicemanager

import (
	"testing"
)

func TestGetMD5Hash(t *testing.T) {
	want := "6980afb078fd03156924b050c1df180a"
	got := md5Hash("my-md5-hash")
	if got != want {
		t.Fatal("could not hash md5")
	}
}
