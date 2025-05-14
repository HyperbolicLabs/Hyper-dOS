package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	// this takes the password from stdin
	cmd := exec.Command("sudo", "-S", "ls")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

}

func alt() {

	// this one works, but it will show the password.
	// TODO how to hide it?
	cmd := exec.Command("sudo", "ls")

	// Create pseudo-terminal
	f, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("Error starting command:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Connect standard input/output to the pty
	go io.Copy(os.Stdout, f)
	go io.Copy(f, os.Stdin)

	// Wait for command to complete
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Command failed:", err)
		os.Exit(1)
	}
}
