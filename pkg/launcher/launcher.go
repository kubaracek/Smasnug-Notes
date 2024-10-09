package launcher

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Launcher interface {
	InstallAppId(string) error
	IsInstalledAppId(string) error
	LaunchAppId(string) (bool, error)
}

type Impl struct{}

func NewLauncher() *Impl {
	return &Impl{}
}

func (_ Impl) InstallAppId(appId string) error {
	// Create the winget command
	cmd := exec.Command("winget", "install", appId, "--accept-package-agreements", "--accept-source-agreements")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// Run the command and capture output/error
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("installation failed: %v", err)
	}

	return nil
}

func (_ Impl) IsInstalledAppId(appId string) (bool, error) {
	cmd := exec.Command("winget", "list")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// Run the command and capture output/error
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to list installed applications: %v", err)
	}

	// Convert the output to a string
	output := string(stdoutStderr)

	// Check if the appId is present in the output
	if strings.Contains(output, appId) {
		return true, nil
	}

	return false, nil
}

func (_ Impl) LaunchAppId(appId string) error {
	samsungNotesCmd := fmt.Sprintf("start shell:AppsFolder\\SAMSUNGELECTRONICSCoLtd.SamsungNotes_%s!App", appId)

	// Execute the command via cmd.exe
	cmd := exec.Command("cmd", "/C", samsungNotesCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launcher failed: %w", err)
	}

	for i := 0; i < 100; i++ {
		processes, err := process.Processes()
		if err != nil {
			return fmt.Errorf("error retrieving processes: %v", err)
		}

		for _, p := range processes {
			name, err := p.Name()
			if err == nil && strings.Contains(strings.ToLower(name), "samsungnotes") {
				fmt.Println("Started OK")
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("app did not start in time")
}
