package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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


func PickImageFile() string {
	return pickFileUbuntu("png|jpg")
}
