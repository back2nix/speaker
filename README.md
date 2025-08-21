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

## ⚙️ Configuration

The application uses a `config.yaml` file for customization. If no config file is found, default values are used.

### Configuration File Structure

Create a `config.yaml` file in the application directory:

```yaml
Speech:
  # Default action when copying text (Ctrl+C)
  # "Translate" - speak translation (default behavior)
  # "Original" - speak original text
  # "RussianOnly" - speak only Russian text
  # Ctrl+Z always performs the opposite action
  DefaultOutput: "RussianOnly"
  En:
    Speed: 3    # English speech speed
    Half: 2     # English speech pitch adjustment
  Ru:
    Speed: 6    # Russian speech speed
    Half: 2     # Russian speech pitch adjustment

Input:
  Device: /dev/input/event1  # Input device for hotkeys
  Listen: ":3111"            # Port for server.go
  Hotkeys:
    Translate: "Ctrl+C"           # Copy and translate
    TranslateOral: "Ctrl+Z"       # Repeat last in opposite language
    ToggleReadMode: "Alt+F"       # Toggle read/translate mode
    TogglePause: "Ctrl+Alt+P"     # Pause/resume
    StopSound: "Alt+C"            # Stop current playback
    ToggleCopyBuffer: "Alt+V"     # Auto-translate Russian → English

Sounds:
  Start: "sound/interface-soft-click-131438.mp3"
  Processing: "sound/computer-processing.mp3"
  Click: "sound/slide-click-92152.mp3"
```

### Configuration Options

#### Speech Settings
- **DefaultOutput**: Controls default behavior when copying text
  - `"Translate"`: Always translate and speak translation
  - `"Original"`: Always speak original text
  - `"RussianOnly"`: Only speak if text is in Russian
- **Speed**: Speech rate (1-10, higher = faster)
- **Half**: Pitch adjustment for voice

#### Input Settings
- **Device**: Path to input device for hotkey detection
- **Listen**: Server listening port
- **Hotkeys**: Customizable key combinations (supports Ctrl, Alt, Shift modifiers)

#### Sound Settings
- **Start**: Sound played when application starts
- **Processing**: Sound during text processing
- **Click**: UI interaction sound

### Default Values
If no config file exists, the application uses these defaults:
- Default output: Translate mode
- English: Speed 3, Half 2
- Russian: Speed 7, Half 2
- Standard hotkey mappings as shown in the table below

## ⌨️ Hotkeys

### 📋 **Clipboard Operations**
| Hotkey | Function |
|:---:|:---|
| `Ctrl+C` | Copy |
| `Ctrl+Shift+C` | Copy and read in English |
| `Alt+C` (×2) | Stop and clear buffer |

### 🔄 **Translation and Repetition**
| Hotkey | Function |
|:---:|:---|
| `Ctrl+Z` | Repeat last in English |
| `Alt+V` | Auto-translate Russian → English |

### ⚙️ **Control**
| Hotkey | Function |
|:---:|:---|
| `Ctrl+F` | Toggle read/translate mode |
| `Ctrl+Alt+P` | Pause |
| `Alt+C` | Stop reading |

### How to check if you're on Wayland
```bash
echo $WAYLAND_DISPLAY
```

### Built with
- https://github.com/nix-community/gomod2nix/blob/master/docs/getting-started.md
