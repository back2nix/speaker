### If Nix is not installed

- https://nixos.org/download/#nix-install-linux
```bash
sh <(curl -L https://nixos.org/nix/install) --daemon
```

### How to run it?

```bash
nix run github:back2nix/speaker
```
- or

```bash
git clone https://github.com/back2nix/speaker
cd speaker
nix run .
```

- or

```bash
git clone https://github.com/back2nix/speaker
cd speaker
nix develop
make run
```

- or

```bash
git clone https://github.com/back2nix/speaker
cd speaker
nix build .
result/bin/speaker
# sudo result/bin/keylogger # run in a separate terminal if you're on Wayland instead of X11
```

#### Voice Translator

- Copy text in a foreign language — the app will read it aloud in Russian.

#### Hotkeys

```
ctrl+c | ctrl+shift+c - copy and read in English
ctrl+z - repeat last in English
ctrl+f - toggle (translate/read only)
ctrl+alt+p - pause
alt+c - stop reading
alt+c x2 - stop and clear clipboard
alt+v - auto-translate copied Russian text to English and put it in the clipboard
```

### How to check if you're on Wayland

```bash
echo $WAYLAND_DISPLAY
```

### Built with

- https://github.com/nix-community/gomod2nix/blob/master/docs/getting-started.md
