use std::sync::mpsc::{self, Receiver};

use crate::job::{JobManager, JobResult};
use crate::scanner::{Action, Scanner};

pub struct SkratoApp {
    scanner: Scanner,
    jobs: JobManager,
    result_rx: Receiver<JobResult>,

    status: String,
    last_result: Option<(bool, Option<i32>)>,
    // race-safety
    run_id: u64,
}

impl SkratoApp {
    pub fn new(_cc: &eframe::CreationContext<'_>) -> Self {
        let (tx, rx) = mpsc::channel();
        let scanner = Scanner::new();
        let jobs = JobManager::new(tx, 0);

        let mut app = Self {
            scanner,
            jobs,
            result_rx: rx,
            status: "Ready.".to_string(),
            last_result: None,
            run_id: 0,
        };

        app.scanner.scan_actions();
        app.status = app.scanner.status.clone();
        app
    }

    fn rescan(&mut self) {
        if self.jobs.running() {
            return;
        }
        self.scanner.scan_actions();
        self.status = self.scanner.status.clone();
    }

    fn start_action(&mut self, action: Action) {
        if self.jobs.running() {
            return;
        }
        self.run_id += 1;
        self.jobs.start(action.clone());
        self.status = format!("{}...", action.label);
        self.last_result = None;
    }

    fn poll_results(&mut self) {
        while let Ok(res) = self.result_rx.try_recv() {
            // race-safety: ignore results from older runs
            if res.run_id != self.run_id {
                continue;
            }

            if res.ok {
                self.status = format!("{} finished successfully.", res.label);
            } else {
                let ec = res
                    .exit_code
                    .map(|v| v.to_string())
                    .unwrap_or("nil".to_string());
                self.status = format!("{} failed (exit code: {}).", res.label, ec);
            }

            self.last_result = Some((res.ok, res.exit_code));
            self.jobs.set_running(false);
        }
    }
}

impl eframe::App for SkratoApp {
    fn update(&mut self, ctx: &eframe::egui::Context, _frame: &mut eframe::Frame) {
        self.poll_results();

        eframe::egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("skrato");
            ui.label("Bootloader and initramfs maintenance");

            ui.separator();

            ui.label(&self.status);

            let last = match self.last_result {
                None => "Last run: (none)".to_string(),
                Some((ok, ec)) => {
                    let ok_text = if ok { "success" } else { "failed" };
                    format!(
                        "Last run: {} (exit code: {})",
                        ok_text,
                        ec.map(|v| v.to_string()).unwrap_or("nil".to_string())
                    )
                }
            };
            ui.label(last);

            ui.separator();
            ui.label("Detected tools:");
            if !self.scanner.detected.is_empty() {
                for d in &self.scanner.detected {
                    ui.label(d);
                }
            } else {
                ui.label("none");
            }

            ui.separator();

            // Actions
            if self.jobs.running() {
                ui.add_enabled(false, eframe::egui::Button::new("Working..."));
            } else {
                for (i, action) in self.scanner.actions.clone().into_iter().enumerate() {
                    let button_text = action.label.clone();
                    let action_clone = action.clone();
                    if ui
                        .add(
                            eframe::egui::Button::new(button_text)
                                .min_size(eframe::egui::vec2(180.0, 38.0)),
                        )
                        .clicked()
                    {
                        // start
                        let _ = i;
                        self.start_action(action_clone);
                        break;
                    }
                }
            }

            if ui
                .add(eframe::egui::Button::new("Rescan").min_size(eframe::egui::vec2(180.0, 38.0)))
                .clicked()
            {
                self.rescan();
            }
        });
    }
}
