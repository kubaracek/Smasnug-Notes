package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Helper function to check if the current user is an admin
func isAdmin() bool {
	token := windows.Token(0)
	return token.IsElevated()
}

// Request elevation using Windows ShellExecute to trigger UAC
func runElevated() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command("cmd", "/C", exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CmdLine:    fmt.Sprintf("/c start runas %s", exe),
	}
	return cmd.Run()
}

// Read a cloak key value
func getRegistryKeyValue(path, name string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err := key.GetStringValue(name)
	return value, err
}

// Set a cloak key value
func setRegistryKeyValue(path, name, value string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.SetStringValue(name, value)
}

// Function to start Samsung Notes using its executable path
func startSamsungNotes() error {
	// Command to start the UWP app using shell:AppsFolder
	samsungNotesCmd := "start shell:AppsFolder\\SAMSUNGELECTRONICSCoLtd.SamsungNotes_wyx1vj98g3asy!App"

	// Execute the command via cmd.exe
	cmd := exec.Command("cmd", "/C", samsungNotesCmd)

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Samsung Notes: %w", err)
	}

	return nil
}

// Wait for Samsung Notes to start by checking process list
func waitForSamsungNotes() error {
	for i := 0; i < 10; i++ {
		processes, err := process.Processes()
		if err != nil {
			return fmt.Errorf("error retrieving processes: %v", err)
		}

		for _, p := range processes {
			name, err := p.Name()
			if err == nil && strings.Contains(strings.ToLower(name), "samsungnotes") {
				fmt.Println("Samsung Notes is running.")
				return nil
			}
		}
		fmt.Println("Waiting for Samsung Notes to start...")
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("Samsung Notes did not start in time")
}

func main() {
	//if !isAdmin() {
	//	fmt.Println("Requesting admin privileges...")
	//	err := runElevated()
	//	if err != nil {
	//		log.Fatalf("Failed to elevate privileges: %v", err)
	//	}
	//	os.Exit(0)
	//}

	var err error

	// Store previous SystemProductName and SystemManufacturer
	prevProductName, err := getRegistryKeyValue(`HARDWARE\DESCRIPTION\System\BIOS`, "SystemProductName")
	if err != nil {
		log.Fatalf("Failed to read SystemProductName: %v", err)
	}
	prevManufacturer, err := getRegistryKeyValue(`HARDWARE\DESCRIPTION\System\BIOS`, "SystemManufacturer")
	if err != nil {
		log.Fatalf("Failed to read SystemManufacturer: %v", err)
	}

	// Ensure the cloak values are reverted even if there's an error
	defer func() {
		fmt.Println("Restoring previous cloak values...")
		err = setRegistryKeyValue(`HARDWARE\DESCRIPTION\System\BIOS`, "SystemProductName", prevProductName)
		if err != nil {
			log.Printf("Failed to restore SystemProductName: %v", err)
		}
		err = setRegistryKeyValue(`HARDWARE\DESCRIPTION\System\BIOS`, "SystemManufacturer", prevManufacturer)
		if err != nil {
			log.Printf("Failed to restore SystemManufacturer: %v", err)
		}
		fmt.Println("Registry values restored.")
	}()

	// Set new SystemProductName and SystemManufacturer
	err = setRegistryKeyValue(`HARDWARE\DESCRIPTION\System\BIOS`, "SystemProductName", "NP960XFG-KC4UK")
	if err != nil {
		log.Fatalf("Failed to set SystemProductName: %v", err)
	}
	err = setRegistryKeyValue(`HARDWARE\DESCRIPTION\System\BIOS`, "SystemManufacturer", "Samsung")
	if err != nil {
		log.Fatalf("Failed to set SystemManufacturer: %v", err)
	}

	// Start Samsung Notes
	err = startSamsungNotes()
	if err != nil {
		log.Fatalf("Failed to start Samsung Notes: %v", err)
	}

	// Wait for Samsung Notes to start
	err = waitForSamsungNotes()
	if err != nil {
		log.Fatalf("Failed to detect Samsung Notes running: %v", err)
	}

	fmt.Println("Process complete. Exiting.")
}
