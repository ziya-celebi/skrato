mod app;
mod job;
mod scanner;

use app::SkratoApp;

fn main() {
    eframe::run_native(
        "skrato",
        eframe::NativeOptions::default(),
        Box::new(|cc| Ok(Box::new(SkratoApp::new(cc)))),
    )
    .unwrap();
}
