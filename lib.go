package main

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/pkg/errors"
	"bytes"
	"log"
	"archive/tar"
	"io"
	"runtime"
)


func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "os error")
	}

	dd := os.Getenv("SNAP_USER_COMMON")

	if strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) || dd == "" {
		dd = filepath.Join(hd, "comments203")
		os.MkdirAll(dd, 0777)
	}

	return dd, nil
}

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func UntestedRandomString(length int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func ExternalLaunch(p string) {
	cmd := "url.dll,FileProtocolHandler"
	runDll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")

	if runtime.GOOS == "windows" {
		exec.Command(runDll32, cmd, p).Run()
	} else if runtime.GOOS == "linux" {
		exec.Command("xdg-open", p).Run()
	}
}

func writeTar(files []string, outPath string) {
	// Create and add some files to the archive.
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, f := range files {
		raw, _ := os.ReadFile(f)

		hdr := &tar.Header{
			Name: f,
			Mode: 0600,
			Size: int64(len(raw)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write([]byte(raw)); err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	os.WriteFile(outPath, buf.Bytes(), 0777)
}

func unpackTar(inPath, outPath string) []string {
	raw, _ := os.ReadFile(inPath)
	buf := bytes.NewBuffer(raw)
	tr := tar.NewReader(buf)

	os.RemoveAll(outPath)
	os.MkdirAll(outPath, 0777)

	ret := make([]string, 0)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}
		fileOutPath := filepath.Join(outPath, filepath.Base(hdr.Name))
		ret = append(ret, filepath.Base(hdr.Name))
		os.MkdirAll(filepath.Dir(fileOutPath), 0777)
		b, err := io.ReadAll(tr)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(fileOutPath, b, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	return ret
}