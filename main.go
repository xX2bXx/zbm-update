package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if os.Geteuid() != 0 {
		log.Fatalf("You must run it as root.\n")
	}

	hasBackup := true
	if getIndexOfBackup() == 0 && len(os.Args) == 3 {
		hasBackup = false
	} else if len(os.Args) != 5 {
		fmt.Printf("usage: \"zbm-update [--target /path/to/ZBM] [--backup /path/to/backup (if you have one)]\"\n")
		return
	}
	
	targetIndex := getIndexOfTarget()
	if targetIndex == 0 {
		log.Fatalf("ERROR: --target argument not found\n")
	}

	targetPath := os.Args[targetIndex]
	if !filepath.IsAbs(targetPath) {
		log.Fatal("ERROR: Target path must be absolute\n")
	}

	if hasBackup {
		backupPath := os.Args[getIndexOfBackup()]
		log.Printf("Creating backup file... ")
		if err := createBackup(targetPath, backupPath); err != nil {
			log.Fatalf("ERROR: %v\n", err)
		}
		log.Printf("Success!\n")
	}

	log.Printf("Installing new ZFSBootMenu...")
	if err := installNewBootloader("https://get.zfsbootmenu.org/efi", targetPath); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
	log.Printf("Success!\nNo errors reported! Enjoy your new version of kernel & bootloader :)")
}

func createBackup(targetPath, backupPath string) error {
	data, err := os.ReadFile(targetPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return err
	}
	return nil
}

func installNewBootloader(url, targetPath string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Download failed: HTTP %d", resp.StatusCode)
	}

	tmpFile := targetPath + ".tmp"
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	if err := os.Chmod(tmpFile, 0644); err != nil {
		os.Remove(tmpFile)
		return err
	}

	if err := os.Rename(tmpFile, targetPath); err != nil {
		os.Remove(tmpFile)
		return err
	}

	return nil
}

func getIndexOfTarget() int {
	for index := range os.Args {
		if os.Args[index] == "--target" && len(os.Args) > index+1 {
			return index + 1
		}
	}
	return 0
}

func getIndexOfBackup() int {
	for index := range os.Args {
		if os.Args[index] == "--backup" && len(os.Args) > index+1 {
			return index + 1
		}
	}
	return 0
}
