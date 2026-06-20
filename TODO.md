# TODO - Make skrato pure Rust (remove LÖVE from installer)

- [ ] Search repo for any remaining references to love/lua
- [ ] Update installer.sh to install/copy the compiled Rust binary instead of requiring love
- [ ] Update README.md accordingly
- [ ] Ensure CI builds release binary (optional but recommended)
- [ ] Test: cargo build --release, then bash installer.sh install, then run skrato

