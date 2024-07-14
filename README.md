### Как запустить?

```bash
nix run github:back2nix/speaker
```
- или

```bash
git clone https://github.com/back2nix/speaker
nix run .
```

- или

```bash
git clone https://github.com/back2nix/speaker
nix develop
go mod tidy
make run
```
- или

```bash
nix build .
result/bin/speaker
# sudo result/bin/keylogger # запустить в отдельном окне если у вас Wayland а не X11
```

#### Переводчик вслух

- Копируешь текст на иностранном языке и программа читает в слух на русском

#### Горячие клавиши

```
ctrl+c | ctrl+shift+c - скопировать и прочитать на английском
ctrl+z - повторить на английском
ctrl+f - переключить (переводчик)/(читать без перевода)
ctrl+alt+p  - Пауза
alt+c - break read
alt+c x2 - break and flush clipboard
```

### Как проверить что у вас Wayland

echo $WAYLAND_DISPLAY

### Собрано с помощью

- https://github.com/nix-community/gomod2nix/blob/master/docs/getting-started.md
