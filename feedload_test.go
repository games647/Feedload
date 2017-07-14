package main

import (
	"testing"
)

func TestExtract(t *testing.T) {
	ext := extractFileExt("http://sub.host.tld/archive/file.mp3")
	if ext != "mp3" {
		t.Error("Expected mp3 got", ext)
	}
}
