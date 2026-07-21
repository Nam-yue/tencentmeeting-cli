//go:build !windows && !openharmony

package keychain

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReadEncFileRejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	realFile := filepath.Join(dir, "real.enc")
	linkFile := filepath.Join(dir, "link.enc")

	_ = os.WriteFile(realFile, []byte("encrypted data"), 0600)
	if err := os.Symlink(realFile, linkFile); err != nil {
		t.Fatalf("os.Symlink() error: %v", err)
	}

	_, err := readEncFile(linkFile)
	if err == nil {
		t.Fatal("readEncFile should reject symlinks")
	}

	data, err := readEncFile(realFile)
	if err != nil {
		t.Fatalf("readEncFile should accept regular file: %v", err)
	}
	if !bytes.Equal(data, []byte("encrypted data")) {
		t.Fatalf("readEncFile data mismatch")
	}
}

func TestReadEncFileUnreadable(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root bypasses permission checks")
	}

	dir := t.TempDir()
	filePath := filepath.Join(dir, "unreadable.enc")
	_ = os.WriteFile(filePath, []byte("encrypted"), 0000)
	defer os.Chmod(filePath, 0600)

	_, err := readEncFile(filePath)
	if err == nil {
		t.Fatal("readEncFile should fail on unreadable file")
	}
}

func TestAtomicWriteFileDirPermission(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "secure-sub")
	filePath := filepath.Join(subDir, "test.enc")

	if err := atomicWriteFile(filePath, []byte("data")); err != nil {
		t.Fatalf("atomicWriteFile() error: %v", err)
	}

	info, err := os.Stat(subDir)
	if err != nil {
		t.Fatalf("os.Stat() error: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != 0700 {
		t.Errorf("directory permission = %o, want 0700", perm)
	}
}
