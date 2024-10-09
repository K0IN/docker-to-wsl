package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/jxeng/shortcut"
	"github.com/yuk7/wsllib-go"
)

func importWsl(distroName string) (err error) {
	_ = wsllib.WslUnregisterDistribution(distroName)
	path, err := filepath.Abs("image.tar")
	if err != nil {
		return err
	}

	var re = regexp.MustCompile(`~[\p{L}0-9\s]+`)
	escapedPath := re.ReplaceAllString(distroName, `-`)

	cmd := exec.Command("wsl", "--import", distroName, fmt.Sprintf("./%s", escapedPath), "image.tar", "--version", "2")
	_, err = cmd.Output()
	if err != nil {
		return wsllib.WslRegisterDistribution(distroName, path)
	}

	return nil
}

func setDefaultWsl(distroName string) (err error) {
	cmd := exec.Command("wsl", "--set-default", distroName)
	_, err = cmd.Output()
	return err
}

func addToStartMenu(distroName string) (err error) {
	appData, exists := os.LookupEnv("APPDATA")
	if !exists {
		return fmt.Errorf("APPDATA environment variable not found")
	}
	programsPath := filepath.Join(appData,
		"Microsoft", "Windows", "Start Menu", "Programs",
		fmt.Sprintf("%s.lnk", distroName))

	sc := shortcut.Shortcut{
		ShortcutPath:     programsPath,
		Target:           "wsl.exe",
		Arguments:        fmt.Sprintf("~ -d %s", distroName),
		IconLocation:     "%SystemRoot%\\System32\\SHELL32.dll,0",
		Description:      "",
		Hotkey:           "",
		WindowStyle:      "1",
		WorkingDirectory: "",
	}
	return shortcut.Create(sc)
}

func launchWsl(distroName string) (err error) {
	_, err = wsllib.WslLaunchInteractive(distroName, "", true)
	return err
}
