package keychain

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func TestAtomicWriteFileConcurrent(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "concurrent.enc")

	const n = 10
	var wg sync.WaitGroup
	errs := make([]error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errs[idx] = atomicWriteFile(filePath, []byte(fmt.Sprintf("writer-%d-payload", idx)))
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("writer %d error: %v", i, err)
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("file should not be empty after concurrent writes")
	}
}

func TestWriteReadEncFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.enc")
	data := []byte("encrypted data content for file roundtrip")

	if err := atomicWriteFile(filePath, data); err != nil {
		t.Fatalf("atomicWriteFile() error: %v", err)
	}

	read, err := readEncFile(filePath)
	if err != nil {
		t.Fatalf("readEncFile() error: %v", err)
	}

	if !bytes.Equal(read, data) {
		t.Fatalf("readEncFile() = %q, want %q", read, data)
	}
}

func TestReadEncFileNotFound(t *testing.T) {
	_, err := readEncFile(filepath.Join(t.TempDir(), "nonexistent.enc"))
	if err != ErrNotFound {
		t.Fatalf("readEncFile() error = %v, want ErrNotFound", err)
	}
}

func TestReadEncFileEmpty(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "empty.enc")
	_ = os.WriteFile(filePath, []byte{}, 0600)

	_, err := readEncFile(filePath)
	if err == nil {
		t.Fatal("readEncFile() should fail on empty file")
	}
}

func TestEncFilePath(t *testing.T) {
	got := encFilePath("/data/tmeet", "sdk123")
	want := filepath.Join("/data/tmeet", "sdk123.enc")
	if got != want {
		t.Fatalf("encFilePath() = %q, want %q", got, want)
	}
}

func TestAtomicWriteFileInvalidDir(t *testing.T) {
	dir := t.TempDir()
	blockFile := filepath.Join(dir, "blocker")
	_ = os.WriteFile(blockFile, []byte("x"), 0600)

	err := atomicWriteFile(filepath.Join(blockFile, "sub", "test.enc"), []byte("data"))
	if err == nil {
		t.Fatal("atomicWriteFile should fail when parent path is a regular file")
	}
}

func TestAtomicWriteFileChmodError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod semantics differ on Windows")
	}
	if os.Getuid() == 0 {
		t.Skip("root bypasses permission checks")
	}

	dir := t.TempDir()
	readonlyDir := filepath.Join(dir, "locked")
	_ = os.MkdirAll(readonlyDir, 0700)
	subDir := filepath.Join(readonlyDir, "sub")
	_ = os.MkdirAll(subDir, 0700)
	_ = os.Chmod(readonlyDir, 0500)
	defer os.Chmod(readonlyDir, 0700)

	_ = atomicWriteFile(filepath.Join(subDir, "test.enc"), []byte("data"))
}

func TestWriteEncFileNoTmpResidue(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "clean.enc")

	_ = atomicWriteFile(filePath, []byte("data"))

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".tmp" {
			t.Errorf("tmp file should be cleaned up: %s", e.Name())
		}
	}
}
