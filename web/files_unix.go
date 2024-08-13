//go:build !windows

package main

import "os/exec"

// TODO change in prod to /var/lib/pva or sth
const directory = "./files"

func setupDirectory() error {
	if err := exec.Command("mkdir", "-p", directory).Run(); err != nil {
		return err
	}

	//if err := exec.Command("chown", "-R", "root", directory).Run(); err != nil {
	//	return err
	//}
	//
	//if err := exec.Command("chmod", "-R", "600", directory).Run(); err != nil {
	//	return err
	//}

	return nil
}
