//go:build !windows

package main

import "os/exec"

const fileDirectory = "/var/lib/pva-server"

func setupDirectory() error {
	if err := exec.Command("mkdir", "-p", fileDirectory).Run(); err != nil {
		return err
	}

	if err := exec.Command("chown", "-R", "root", fileDirectory).Run(); err != nil {
		return err
	}

	if err := exec.Command("chmod", "-R", "600", fileDirectory).Run(); err != nil {
		return err
	}

	return nil
}
