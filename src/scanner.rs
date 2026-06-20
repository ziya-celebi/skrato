use std::process::Command;

#[derive(Clone, Debug)]
pub struct Action {
    pub label: String,
    pub command: String,
}

pub struct Scanner {
    pub actions: Vec<Action>,
    pub detected: Vec<String>,
    pub status: String,
}

impl Scanner {
    pub fn new() -> Self {
        Self {
            actions: Vec::new(),
            detected: Vec::new(),
            status: "Ready.".to_string(),
        }
    }

    pub fn scan_actions(&mut self) {
        self.actions.clear();
        self.detected.clear();

        if !sh_success("command -v pkexec >/dev/null 2>&1") {
            self.status = "pkexec not found. Install polkit first.".to_string();
            return;
        }

        let grub = if sh_success("command -v grub-mkconfig >/dev/null 2>&1") {
            Some("grub-mkconfig")
        } else if sh_success("command -v grub2-mkconfig >/dev/null 2>&1") {
            Some("grub2-mkconfig")
        } else {
            None
        };

        if let Some(grub) = grub {
            if sh_success("[ -e '/boot/grub' ]") {
                self.add_action(
                    "Update GRUB",
                    &build_pkexec_command(grub, &["-o", "/boot/grub/grub.cfg"]),
                    "GRUB",
                );
            } else if sh_success("[ -e '/boot/grub2' ]") {
                self.add_action(
                    "Update GRUB",
                    &build_pkexec_command(grub, &["-o", "/boot/grub2/grub.cfg"]),
                    "GRUB",
                );
            }
        }

        if sh_success("command -v bootctl >/dev/null 2>&1")
            && (sh_success("[ -e '/boot/loader' ]") || sh_success("[ -e '/efi/loader' ]"))
        {
            self.add_action(
                "Update systemd-boot",
                &build_pkexec_command("bootctl", &["update"]),
                "systemd-boot",
            );
        }

        if sh_success("command -v refind-install >/dev/null 2>&1") {
            self.add_action(
                "Refresh rEFInd",
                &build_pkexec_command("refind-install", &[]),
                "rEFInd",
            );
        }

        if sh_success("command -v mkinitcpio >/dev/null 2>&1") {
            self.add_action(
                "Rebuild initramfs",
                &build_pkexec_command("mkinitcpio", &["-P"]),
                "mkinitcpio",
            );
        } else if sh_success("command -v dracut >/dev/null 2>&1") {
            self.add_action(
                "Rebuild initramfs",
                &build_pkexec_command("dracut", &["--regenerate-all", "--force"]),
                "dracut",
            );
        } else if sh_success("command -v update-initramfs >/dev/null 2>&1") {
            self.add_action(
                "Rebuild initramfs",
                &build_pkexec_command("update-initramfs", &["-u", "-k", "all"]),
                "update-initramfs",
            );
        }

        self.status = if !self.actions.is_empty() {
            format!("Detected: {}", self.detected.join(", "))
        } else {
            "Nothing supported detected.".to_string()
        };
    }

    fn add_action(&mut self, label: &str, command: &str, detect_name: &str) {
        self.actions.push(Action {
            label: label.to_string(),
            command: command.to_string(),
        });
        self.detected.push(detect_name.to_string());
    }
}

fn sh_success(command: &str) -> bool {
    // We use /bin/sh via shell -c so redirects/[] checks work.
    Command::new("sh")
        .arg("-c")
        .arg(command)
        .status()
        .map(|s| s.success())
        .unwrap_or(false)
}

fn quote_single_posix(s: &str) -> String {
    // safe quoting for POSIX sh single quotes
    // 'foo' => '\'' pattern
    format!("'{}'", s.replace("'", "'\\''"))
}

fn build_pkexec_command(binary: &str, args: &[&str]) -> String {
    let mut cmd = format!("pkexec {}", quote_single_posix(binary));
    for a in args {
        cmd.push(' ');
        cmd.push_str(&quote_single_posix(a));
    }
    cmd
}
