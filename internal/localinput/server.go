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

func Start(cancel context.CancelFunc, translator intf.Translator) (err error) {
	readRU := false
	flagToCopyBuffer := false

	translateshell.Play(findSound("sound/interface-soft-click-131438.mp3"))

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
			translateshell.Play(findSound("sound/computer-processing.mp3"))
		} else {
			translateshell.Play(findSound("sound/slide-click-92152.mp3"))
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
			translateshell.Play(findSound("sound/computer-processing.mp3"))
		} else {
			translateshell.Play(findSound("sound/slide-click-92152.mp3"))
		}
	}

	ctrl_p_func = func() {
		if !translator.CheckPause() {
			translateshell.Play(findSound("sound/slide-click-92152.mp3"))
		} else {
			translateshell.Play(findSound("sound/slide-click-92152.mp3"))
		}
		translator.SetPause()
	}

	return devInput()
}

func findSound(filename string) string {
	// Сначала проверяем текущую директорию
	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	// Затем проверяем в поддиректории sound текущей директории
	soundInCurrent := filepath.Join("sound", filename)
	if _, err := os.Stat(soundInCurrent); err == nil {
		return soundInCurrent
	}

	// Если не найдено, ищем в PATH
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, path := range paths {
		fullPath := filepath.Join(path, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	// Если не найдено нигде, возвращаем исходное имя файла
	return filename
}

var channel = make(chan map[uint16]bool)

// sudo chmod 666 /dev/input/event0
func devInput() (err error) {
	dev, err := evdev.Open("/dev/input/event0")
	if err != nil {
		return err
	}
	defer dev.File.Close()

	go PresedWorker()

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
					// fmt.Println("Key_Down: %d", keyCode)
					key[keyCode] = true
					// fmt.Println(key)
					channel <- key
				}
				if event.Value == 0 {
					// fmt.Println("Key_UP: %d", keyCode)
					if _, ok := key[keyCode]; ok {
						delete(key, keyCode)
						// fmt.Println(key)
						channel <- key
					}
				}
			}
		}
	}
}

func PresedWorker() {
	var last string
	for k := range channel {
		var mess string
		switch {
		case CheckKeys(evdev.KEY_P, true, false, true, k):
			mess = "L_CTRL+L_ALT+P"
		case CheckKeys(evdev.KEY_C, true, true, false, k):
			mess = "L_CTRL+L_SHIFT+C"
		case CheckKeys(evdev.KEY_C, true, false, false, k):
			mess = "L_CTRL+C"
		case CheckKeys(evdev.KEY_Z, true, false, false, k):
			mess = "L_CTRL+Z"
		case CheckKeys(evdev.KEY_Z, false, false, true, k):
			mess = "L_CTRL+Z"
		case CheckKeys(evdev.KEY_F, false, false, true, k):
			mess = "L_ALT+F"
		case CheckKeys(evdev.KEY_C, false, false, true, k):
			if last == "L_ALT+C" {
				mess = "L_ALT+Cx2"
				last = ""
			} else {
				mess = "L_ALT+C"
			}
		case CheckKeys(evdev.KEY_V, false, false, true, k):
			mess = "L_ALT+V"
		}
		if mess != "" {
			last = mess
			SendMessage(mess)
		}
	}
}

func SendMessage(input string) {
	if input == "" {
		return
	}

	switch input {
	case "L_CTRL+L_SHIFT+C", "L_CTRL+C":
		ctrl_c_func()
	case "L_CTRL+Z", "L_ALT+Z":
		ctrl_z_func()
	case "L_ALT+F":
		ctrl_f_func()
	case "L_CTRL+L_ALT+P":
		ctrl_p_func()
	case "L_ALT+C":
		alt_c_func(0)
	case "L_ALT+V":
		alt_v_func()
	case "L_ALT+Cx2":
		alt_c_func(1)
	}
}

func CheckKeys(key uint16, ctrl, shift, alt bool, k map[uint16]bool) bool {
	var ctrlL, altL, shiftL bool
	_, ctrlL = k[evdev.KEY_LEFTCTRL]
	_, altL = k[evdev.KEY_LEFTALT]
	_, shiftL = k[evdev.KEY_LEFTSHIFT]
	if _, ok := k[key]; ok && ctrl == ctrlL && shift == shiftL && alt == altL {
		return true
	}
	return false
}

var (
	reg0             = regexp.MustCompile(`[^a-zA-Z\p{Han}0-9 .,\r\n]+`)
	reg2             = regexp.MustCompile(`([\p{L}])\.([\p{L}])`)
	reg3             = regexp.MustCompile(`([[:lower:]])([[:upper:]])`)
	reg4             = regexp.MustCompile(`(\b(\p{L}+)\b)`)
	singleSpaceRegex = regexp.MustCompile(`\s+`)

	re = regexp.MustCompile(`([\wλ]+)\s+∈\s+([\w ]+)`)

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
		// "ei":   "Exponential constant e raised to the power i",
		// "dim":  "Dimension",
		"rk(A)":    "Rank of matrix A",
		"Im(Φ)":    "Image of transformation Φ",
		"ker(Φ)":   "Kernel of transformation Φ",
		"span[b1]": "Span of vector b1",
		"tr(A)":    "Trace of matrix A",
		"det(A)":   "Determinant of matrix A",
		// "| · |": "Absolute value",
		// "∥·∥": "Norm",
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
	var regexPattern string
	for symbol := range mathSymbols {
		if regexPattern != "" {
			regexPattern += "|"
		}
		regexPattern += regexp.QuoteMeta(symbol)
	}
	regexPattern = "(" + regexPattern + ")"

	// Заменяем символы в строке на их описания
	return regexp.MustCompile(regexPattern)
}

func MathRegex(input string) (out string, err error) {
	// out = re.ReplaceAllStringFunc(input, func(match string) string {
	// 	matches := re.FindStringSubmatch(match)
	// 	element := matches[1]
	// 	set := matches[2]
	//
	// 	if replacement, exists := dictionary[element]; exists {
	// 		return fmt.Sprintf("%s belongs to the set of %s", replacement, set)
	// 	}
	// 	return match
	// })

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
