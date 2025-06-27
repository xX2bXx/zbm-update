package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if os.Geteuid() != 0 {
		log.Fatalf("You must run it as root.\n")
	}

	hasBackup := true
	if getIndexOf("--backup") == 0 && (len(os.Args) == 3 || len(os.Args) == 5) {
		hasBackup = false
	} else if len(os.Args) != 7 {
		fmt.Printf("usage: \"zbm-update --target /path/to/ZBM.EFI [--backup /path/to/BACKUP.EFI (if you have one)] [--fallback true|false]\"\n")
		return
	}

	targetPath := os.Args[getIndexOf("--target")]
	if getIndexOf("--target") == 0 {
		log.Fatalf("ERROR: --target argument not found\n")
	}

	if !filepath.IsAbs(targetPath) {
		log.Fatal("ERROR: Target path must be absolute\n")
	}

	backupPath := os.Args[getIndexOf("--backup")]
	if hasBackup {
		log.Printf("Creating backup file... ")
		if err := createBackup(targetPath, backupPath); err != nil {
			log.Fatalf("ERROR: %v\n", err)
		}
		log.Printf("Success!\n")
	}

	log.Printf("Installing new ZFSBootMenu...")
	if err := installNewBootloader(targetPath); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	log.Printf("Success!\n")

	if value, err := strconv.ParseBool(os.Args[getIndexOf("--fallback")]); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	} else if value {
		log.Printf("Creating fallback (/boot/efi/EFI/BOOT/BOOTX64.EFI)...")
		if err := createFallback(hasBackup, targetPath, backupPath); err != nil {
			log.Fatalf("ERROR: %v\n", err)
		}	
		log.Printf("Success!\n")
	}
	log.Printf("No errors reported! Enjoy your updated ZFSBootMenu :)\n")
}

func createBackup(targetPath, backupPath string) error {
	data, err := os.ReadFile(targetPath)
	if err != nil {
		return err
	}
	
	return os.WriteFile(backupPath, data, 0600)
}

func installNewBootloader(targetPath string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get("https://get.zfsbootmenu.org/efi")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Download failed: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	tmpFile := targetPath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmpFile, targetPath); err != nil {
		os.Remove(tmpFile)
		return err
	}

	return nil
}

func createFallback(hasBackup bool, targetPath, backupPath string) error {
	var data []byte
	var err error
	
	if hasBackup {
		data, err = os.ReadFile(backupPath)
	} else {
		data, err = os.ReadFile(targetPath)
	}
	if err != nil {
		return err
	}

	if err := os.MkdirAll("/boot/efi/EFI/BOOT", 0755); err != nil {
		return err
	}

	return os.WriteFile("/boot/efi/EFI/BOOT/BOOTX64.EFI", data, 0644)
}

func getIndexOf(flag string) int {
	for index := range os.Args {
		if os.Args[index] == flag && len(os.Args) > index+1 {
			return index + 1
		}
	}
	return 0
}
