package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Matyjash/DirSync/dirsync"
)

func main() {
	fmt.Println("DirSync")

	deleteMissing, source, destination := parseFlags()

	if err := validateDirPath(source); err != nil {
		fmt.Printf("invalid source path: %v", err)
		os.Exit(1)
	}
	if err := validateDirPath(destination); err != nil {
		fmt.Printf("invalid destination path: %v", err)
		os.Exit(1)
	}

	dirSync := dirsync.NewDirSync(deleteMissing)
	dirSync.Sync(source, destination)

	fmt.Println("synchronization completed")
}

func parseFlags() (deleteMissing bool, source, destination string) {
	deleteMissingFlag := flag.Bool("delete-missing", false, "Delete files in destination that are not present in source")
	flag.Parse()

	fmt.Printf("Delete missing files: %v\n", *deleteMissingFlag)

	args := flag.Args()
	if len(args) < 2 || len(args) > 3 {
		fmt.Printf("%d", len(args))
		fmt.Printf("Usage: %s <--delete-missing> <source> <destination> ", os.Args[0])
		os.Exit(1)
	}

	return *deleteMissingFlag, args[0], args[1]
}

func validateDirPath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}
	return nil
}
