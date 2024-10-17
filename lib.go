package main

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"fmt"

	// "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pkg/errors"
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


func GetPickerPath() string {
	homeDir, _ := os.UserHomeDir()
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = filepath.Join(homeDir, "bin", "fpicker")
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", "fpicker")
	}

	return cmdPath
}


func pickFileUbuntu(exts string) string {
	fPickerPath := GetPickerPath()

	rootPath, _ := GetRootPath()
	cmd := exec.Command(fPickerPath, rootPath, exts)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}
