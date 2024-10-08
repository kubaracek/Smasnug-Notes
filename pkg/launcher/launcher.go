package launcher

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os/exec"
	"strings"
	"time"
)

type Launcher interface {
	InstallApp() error
	LaunchApp() error
}

const SAMSUNG_NOTES_INSTALL_ID = "9nblggh43vhv"
const SAMSUNG_NOTES_LAUNCH_ID = "wyx1vj98g3asy"

type SamsungNotes struct{}

func NewSamsungNotes() *SamsungNotes {
	return &SamsungNotes{}
}

func (_ SamsungNotes) InstallApp() error {
	// Create the winget command
	cmd := exec.Command("winget", "install", SAMSUNG_NOTES_INSTALL_ID, "--accept-package-agreements", "--accept-source-agreements")

	// Run the command and capture output/error
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("installation failed: %v", err)
	}

	return nil
}

func (_ SamsungNotes) LaunchApp() error {
	samsungNotesCmd := fmt.Sprintf("start shell:AppsFolder\\SAMSUNGELECTRONICSCoLtd.SamsungNotes_%s!App", SAMSUNG_NOTES_LAUNCH_ID)

	fmt.Println("Launching Samsung Notes")
	// Execute the command via cmd.exe
	cmd := exec.Command("cmd", "/C", samsungNotesCmd)

	fmt.Println("Waiting for process to start...")

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
