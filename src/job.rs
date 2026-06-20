use std::process::Command;
use std::thread;

use crate::scanner::Action;

#[derive(Debug, Clone)]
pub struct JobResult {
    pub run_id: u64,
    pub ok: bool,
    pub exit_code: Option<i32>,
    pub label: String,
}

pub struct JobManager {
    pub running: bool,
    run_id: u64,
    result_tx: std::sync::mpsc::Sender<JobResult>,
}

impl JobManager {
    pub fn new(result_tx: std::sync::mpsc::Sender<JobResult>, run_id_start: u64) -> Self {
        Self {
            running: false,
            run_id: run_id_start,
            result_tx,
        }
    }

    pub fn start(
        &mut self,
        action: Action,
        // sender side is stored in self
    ) {
        if self.running {
            return;
        }
        self.running = true;
        self.run_id += 1;
        let rid = self.run_id;
        let tx = self.result_tx.clone();

        thread::spawn(move || {
            // Run in shell so pkexec string is interpreted just like in the Lua version.
            let status = Command::new("sh").arg("-c").arg(action.command).status();

            let (ok, exit_code) = match status {
                Ok(s) => (s.success(), s.code()),
                Err(_) => (false, None),
            };

            let _ = tx.send(JobResult {
                run_id: rid,
                ok,
                exit_code,
                label: action.label,
            });
        });
    }

    pub fn running(&self) -> bool {
        self.running
    }

    pub fn set_running(&mut self, v: bool) {
        self.running = v;
    }
}
