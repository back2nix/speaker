package localinput

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	evdev "github.com/back2nix/golang-evdev"
	"github.com/sirupsen/logrus"

	"github.com/back2nix/speaker/internal/config"
	"github.com/back2nix/speaker/internal/intf"
	"github.com/back2nix/speaker/internal/translateshell"
)

var (
	ctrl_c_func func()
	ctrl_z_func func()
	ctrl_f_func func()
	alt_c_func  func(a int)
	alt_v_func  func()
	ctrl_p_func func()
)

func Start(cancel context.CancelFunc, translator intf.Translator, cfg *config.Config) (err error) {
	readRU := false
	flagToCopyBuffer := false

	logrus.Debug("Calling Start function in localinput")
	// ИСПОЛЬЗУЕМ КЛЮЧИ В НИЖНЕМ РЕГИСТРЕ
	startSoundPath := findSound(cfg.Sounds["start"])
	logrus.WithField("path", startSoundPath).Info("Playing start sound")
	translateshell.Play(startSoundPath)

	ctrl_c_func = func() {
		if translator.CheckPause() {
			return
		}

		time.Sleep(time.Millisecond * 50)
		text, err := clipboard.ReadAll()
		fmt.Println("text:", text)
		if text == "" {
			return
		}

		if flagToCopyBuffer {
			engTxt := translator.Speak(context.Background(), text, `trans -b -t en "%s"`)
			engTxt = strings.TrimSuffix(engTxt, "\n")
			fmt.Printf("eng: %v\n", engTxt)
			if engTxt != "" {
				clipboard.WriteAll(engTxt)
			}
		}

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("clipboard")

			translator.OnlyOriginalRu()
			translator.Go("не скопировалось")
			return
		}

		processedString, err := RegexWork(text)

		fmt.Println(processedString)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("regexp")
			return
		}
		if readRU {
			processedString, _ := RegexWorkRu(text)
			translator.OnlyOriginalRu()
			translator.Go(processedString)
		} else {
			translator.OnlyTranslate()
			translator.Go(processedString)
		}
	}

	ctrl_z_func = func() {
		if translator.CheckPause() {
			return
		}

		time.Sleep(time.Millisecond * 50)
		text, err := clipboard.ReadAll()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("clipboard")

			translator.OnlyOriginalRu()
			translator.Go("не скопировалось")
			return
		}

		processedString, err := RegexWork(text)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("regexp")
			return
		}
		translator.OnlyOriginal()
		translator.Go(processedString)
	}

	ctrl_f_func = func() {
		if translator.CheckPause() {
			return
		}
		readRU = !readRU

		if readRU {
			// ИСПОЛЬЗУЕМ КЛЮЧИ В НИЖНЕМ РЕГИСТРЕ
			soundPath := findSound(cfg.Sounds["processing"])
			logrus.WithField("path", soundPath).Info("Playing processing sound")
			translateshell.Play(soundPath)
		} else {
			// ИСПОЛЬЗУЕМ КЛЮЧИ В НИЖНЕМ РЕГИСТРЕ
			soundPath := findSound(cfg.Sounds["click"])
			logrus.WithField("path", soundPath).Info("Playing click sound")
			translateshell.Play(soundPath)
		}
	}

	alt_c_func = func(a int) {
		translateshell.Stop()
		switch a {
		case 1:
			clipboard.WriteAll("")
		}
	}
	alt_v_func = func() {
		flagToCopyBuffer = !flagToCopyBuffer
		if flagToCopyBuffer {
			// ИСПОЛЬЗУЕМ КЛЮЧИ В НИЖНЕМ РЕГИСТРЕ
			soundPath := findSound(cfg.Sounds["processing"])
			logrus.WithField("path", soundPath).Info("Playing processing sound for copy buffer toggle")
			translateshell.Play(soundPath)
		} else {
			// ИСПОЛЬЗУЕМ КЛЮЧИ В НИЖНЕМ РЕГИСТРЕ
			soundPath := findSound(cfg.Sounds["click"])
			logrus.WithField("path", soundPath).Info("Playing click sound for copy buffer toggle")
			translateshell.Play(soundPath)
		}
	}

	ctrl_p_func = func() {
		// ИСПОЛЬЗУЕМ КЛЮЧИ В НИЖНЕМ РЕГИСТРЕ
		soundPath := findSound(cfg.Sounds["click"])
		if !translator.CheckPause() {
			logrus.WithField("path", soundPath).Info("Playing click sound for pause")
			translateshell.Play(soundPath)
		} else {
			logrus.WithField("path", soundPath).Info("Playing click sound for unpause")
			translateshell.Play(soundPath)
		}
		translator.SetPause()
	}

	return devInput(cfg)
}

func findSound(filename string) string {
	log := logrus.WithField("original_filename", filename)
	log.Debug("Attempting to find sound file")

	if filename == "" {
		log.Warn("Received empty filename for sound")
		return ""
	}

	if _, err := os.Stat(filename); err == nil {
		log.WithField("path", filename).Debug("File found at exact path")
		return filename
	}

	if !strings.ContainsAny(filename, "/\\") {
		soundInCurrent := filepath.Join("sound", filename)
		if _, err := os.Stat(soundInCurrent); err == nil {
			log.WithField("path", soundInCurrent).Debug("File found in 'sound/' directory")
			return soundInCurrent
		}
	}

	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, path := range paths {
		fullPath := filepath.Join(path, filename)
		if _, err := os.Stat(fullPath); err == nil {
			log.WithField("path", fullPath).Debug("File found in system PATH")
			return fullPath
		}
	}

	log.Warn("Sound file not found in any checked location")
	return filename
}

var channel = make(chan map[uint16]bool)

func devInput(cfg *config.Config) (err error) {
	dev, err := evdev.Open(cfg.Input.Device)
	if err != nil {
		return err
	}
	defer dev.File.Close()

	go PresedWorker(cfg)

	key := make(map[uint16]bool)

	for {
		events, err := dev.Read()
		if err != nil {
			panic(err)
		}
		for i := range events {
			event := &events[i]
			if event.Type == evdev.EV_KEY {
				keyCode := event.Code
				if event.Value == 1 {
					key[keyCode] = true
					channel <- key
				}
				if event.Value == 0 {
					if _, ok := key[keyCode]; ok {
						delete(key, keyCode)
						channel <- key
					}
				}
			}
		}
	}
}

func isHotkeyPressed(hotkeyString string, pressedKeys map[uint16]bool) bool {
	parts := strings.Split(hotkeyString, "+")
	if len(parts) == 0 {
		return false
	}

	var mainKey string
	wantsCtrl := false
	wantsAlt := false
	wantsShift := false

	for _, part := range parts {
		p := strings.TrimSpace(strings.ToUpper(part))
		switch p {
		case "CTRL":
			wantsCtrl = true
		case "ALT":
			wantsAlt = true
		case "SHIFT":
			wantsShift = true
		default:
			mainKey = p
		}
	}

	if mainKey == "" {
		return false
	}

	keyCode, ok := config.KeyToCode[mainKey]
	if !ok {
		return false
	}

	if _, keyIsPressed := pressedKeys[keyCode]; !keyIsPressed {
		return false
	}

	_, ctrlLPressed := pressedKeys[evdev.KEY_LEFTCTRL]
	_, ctrlRPressed := pressedKeys[evdev.KEY_RIGHTCTRL]
	isCtrlPressed := ctrlLPressed || ctrlRPressed

	_, altLPressed := pressedKeys[evdev.KEY_LEFTALT]
	_, altRPressed := pressedKeys[evdev.KEY_RIGHTALT]
	isAltPressed := altLPressed || altRPressed

	_, shiftLPressed := pressedKeys[evdev.KEY_LEFTSHIFT]
	_, shiftRPressed := pressedKeys[evdev.KEY_RIGHTSHIFT]
	isShiftPressed := shiftLPressed || shiftRPressed

	return wantsCtrl == isCtrlPressed && wantsAlt == isAltPressed && wantsShift == isShiftPressed
}

func PresedWorker(cfg *config.Config) {
	var lastFiredAction string
	var lastFiredTime time.Time
	sentActions := make(map[string]bool)

	for k := range channel {
		if len(k) == 0 {
			sentActions = make(map[string]bool)
			continue
		}

		var currentAction string
		if isHotkeyPressed(cfg.Input.Hotkeys.TogglePause, k) {
			currentAction = "TogglePause"
		} else if isHotkeyPressed(cfg.Input.Hotkeys.Translate, k) {
			currentAction = "Translate"
		} else if isHotkeyPressed(cfg.Input.Hotkeys.TranslateOral, k) {
			currentAction = "TranslateOral"
		} else if isHotkeyPressed(cfg.Input.Hotkeys.ToggleReadMode, k) {
			currentAction = "ToggleReadMode"
		} else if isHotkeyPressed(cfg.Input.Hotkeys.ToggleCopyBuffer, k) {
			currentAction = "ToggleCopyBuffer"
		} else if isHotkeyPressed(cfg.Input.Hotkeys.StopSound, k) {
			if lastFiredAction == "StopSound" && time.Since(lastFiredTime) < 500*time.Millisecond {
				currentAction = "ClearClipboard"
			} else {
				currentAction = "StopSound"
			}
		}

		if currentAction != "" {
			if !sentActions[currentAction] {
				SendMessage(currentAction)
				sentActions[currentAction] = true
				lastFiredAction = currentAction
				lastFiredTime = time.Now()
			}
		}
	}
}

func SendMessage(action string) {
	if action == "" {
		return
	}

	switch action {
	case "Translate":
		ctrl_c_func()
	case "TranslateOral":
		ctrl_z_func()
	case "ToggleReadMode":
		ctrl_f_func()
	case "TogglePause":
		ctrl_p_func()
	case "StopSound":
		alt_c_func(0)
	case "ClearClipboard":
		alt_c_func(1)
	case "ToggleCopyBuffer":
		alt_v_func()
	}
}

// --- ВОССТАНОВЛЕННЫЙ КОД ---

var (
	reg0             = regexp.MustCompile(`[^a-zA-Z\p{Han}0-9 .,\r\n]+`)
	reg2             = regexp.MustCompile(`([\p{L}])\.([\p{L}])`)
	reg3             = regexp.MustCompile(`([[:lower:]])([[:upper:]])`)
	reg4             = regexp.MustCompile(`(\b(\p{L}+)\b)`)
	singleSpaceRegex = regexp.MustCompile(`\s+`)

	mathSymbols = map[string]string{
		"ё":    "е",
		"∂f":   "Partial derivative of f",
		"∂x":   "Partial derivative of x",
		"⊥":    "Perpendicular",
		"θ":    "Theta",
		"∇":    "Gradient",
		"λ":    "Lambda",
		"Im":   "Imaginary part",
		"0m,n": "Zero matrix with m rows and n columns",
		"0m":   "Zero matrix",
		"1m,n": "Unit matrix with m rows and n columns",
		"1m":   "Unit matrix",
		"rk(A)":    "Rank of matrix A",
		"Im(Φ)":    "Image of transformation Φ",
		"ker(Φ)":   "Kernel of transformation Φ",
		"span[b1]": "Span of vector b1",
		"tr(A)":    "Trace of matrix A",
		"det(A)":   "Determinant of matrix A",
		"Eλ":     "Eigenvalue matrix E for eigenvalue λ",
		"∅":      "Empty set",
		"a ∈ A":  "a belongs to set A",
		"∈":      "belongs to set",
		"α":      "Alpha",
		"β":      "Beta",
		"γ":      "Gamma",
		"𝑥":      "x",
		"𝑦":      "y",
		"𝑧":      "z",
		"𝐴":      "A",
		"𝐵":      "B",
		"𝐶":      "C",
		"𝑥ᵀ":     "Transpose of vector x",
		"𝐴ᵀ":     "Transpose of matrix A",
		"𝐴⁻¹":    "Inverse of matrix A",
		"⟨𝑥, 𝑦⟩": "Inner product of x and y",
		"𝑥ᵀ𝑦":    "Dot product of x and y",
		"ℤ":      "Integers",
		"ℕ":      "Natural numbers",
		"ℝ":      "Real numbers",
		"ℂ":      "Complex numbers",
		"ℝⁿ":     "n-dimensional vector space of real numbers",
		"∈ R ":   "belongs to Real numbers",
		"∈ Rn":   "belongs to n-dimensional vector space of real numbers",
		"∀x":     "For all x",
		"∃x":     "Exists x",
		"a := b": "Assignment a := b",
		"a =: b": "Assignment a =: b",
		"a ∝ b":  "a is proportional to b",
		"g ◦ f":  "Composition of functions g and f",
		"⇐⇒":     "If and only if",
	}

	regexPattern = compileMathSymbolsRegexp()
)

func compileMathSymbolsRegexp() *regexp.Regexp {
	var pattern string
	for symbol := range mathSymbols {
		if pattern != "" {
			pattern += "|"
		}
		pattern += regexp.QuoteMeta(symbol)
	}
	pattern = "(" + pattern + ")"
	return regexp.MustCompile(pattern)
}

func MathRegex(input string) (out string, err error) {
	out = regexPattern.ReplaceAllStringFunc(input, func(matched string) string {
		return mathSymbols[matched]
	})
	return out, nil
}

func RegexWork(tt string) (out string, err error) {
	tt, _ = MathRegex(tt)
	tt = strings.NewReplacer("\n", "", "\r", "", "\"", "", "'", "").Replace(tt)
	tt = reg0.ReplaceAllString(tt, " ")
	tt = reg4.ReplaceAllString(tt, " $1 ")
	tt = reg3.ReplaceAllString(tt, "$1 $2")
	tt = reg2.ReplaceAllString(tt, "$1. $2")
	tt = strings.NewReplacer(" .", ".", " ,", ",").Replace(tt)
	tt = singleSpaceRegex.ReplaceAllString(tt, " ")
	return strings.TrimSpace(tt), nil
}

var (
	reg01              = regexp.MustCompile(`[^а-яА-Яa-zA-Z0-9 .,]+`)
	singleSpacePattern = regexp.MustCompile(`\s+`)
)

func RegexWorkRu(tt string) (out string, err error) {
	tt = strings.ReplaceAll(tt, "ё", "е")
	tt = reg01.ReplaceAllString(tt, " ")
	tt = singleSpacePattern.ReplaceAllString(tt, " ")
	tt = strings.ReplaceAll(tt, " .", ".")
	tt = strings.ReplaceAll(tt, " ,", ",")
	return strings.TrimSpace(tt), nil
}
