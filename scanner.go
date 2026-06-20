package main

import (
	"os/exec"
	"strings"
)

type Action struct {
	Label   string
	Command string
}

type Scanner struct {
	Actions  []Action
	Detected []string
	Status   string
}

func NewScanner() *Scanner {
	return &Scanner{
		Actions:  []Action{},
		Detected: []string{},
		Status:   "Ready.",
	}
}

func (s *Scanner) ScanActions() {
	s.Actions = nil
	s.Detected = nil

	if !shSuccess("command -v pkexec >/dev/null 2>&1") {
		s.Status = "pkexec not found. Install polkit first."
		return
	}

	grub := ""
	if shSuccess("command -v grub-mkconfig >/dev/null 2>&1") {
		grub = "grub-mkconfig"
	} else if shSuccess("command -v grub2-mkconfig >/dev/null 2>&1") {
		grub = "grub2-mkconfig"
	}

	if grub != "" {
		if shSuccess("[ -e '/boot/grub' ]") {
			s.addAction("Update GRUB", buildPkexecCommand(grub, []string{"-o", "/boot/grub/grub.cfg"}), "GRUB")
		} else if shSuccess("[ -e '/boot/grub2' ]") {
			s.addAction("Update GRUB", buildPkexecCommand(grub, []string{"-o", "/boot/grub2/grub.cfg"}), "GRUB")
		}
	}

	if shSuccess("command -v bootctl >/dev/null 2>&1") &&
		(shSuccess("[ -e '/boot/loader' ]") || shSuccess("[ -e '/efi/loader' ]")) {
		s.addAction("Update systemd-boot", buildPkexecCommand("bootctl", []string{"update"}), "systemd-boot")
	}

	if shSuccess("command -v refind-install >/dev/null 2>&1") {
		s.addAction("Refresh rEFInd", buildPkexecCommand("refind-install", []string{}), "rEFInd")
	}

	if shSuccess("command -v mkinitcpio >/dev/null 2>&1") {
		s.addAction("Rebuild initramfs", buildPkexecCommand("mkinitcpio", []string{"-P"}), "mkinitcpio")
	} else if shSuccess("command -v dracut >/dev/null 2>&1") {
		s.addAction("Rebuild initramfs", buildPkexecCommand("dracut", []string{"--regenerate-all", "--force"}), "dracut")
	} else if shSuccess("command -v update-initramfs >/dev/null 2>&1") {
		s.addAction("Rebuild initramfs", buildPkexecCommand("update-initramfs", []string{"-u", "-k", "all"}), "update-initramfs")
	}

	if len(s.Actions) > 0 {
		s.Status = "Detected: " + strings.Join(s.Detected, ", ")
	} else {
		s.Status = "Nothing supported detected."
	}
}

func (s *Scanner) addAction(label string, command string, detectName string) {
	s.Actions = append(s.Actions, Action{
		Label:   label,
		Command: command,
	})
	s.Detected = append(s.Detected, detectName)
}

func shSuccess(command string) bool {
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return cmd.ProcessState.Success()
}

func quoteSinglePosix(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func buildPkexecCommand(binary string, args []string) string {
	cmd := "pkexec " + quoteSinglePosix(binary)
	for _, a := range args {
		cmd += " " + quoteSinglePosix(a)
	}
	return cmd
}