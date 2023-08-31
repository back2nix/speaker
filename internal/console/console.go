package console

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	// "github.com/eiannone/keyboard"

	"github.com/REPO_DEPRECATED/speaker_alpine/internal/intf"
	hook "github.com/robotn/gohook"
	"github.com/sirupsen/logrus"
)

type model struct {
	LogLevel    int
	MaxLogLevel int
}

var mod = model{
	LogLevel:    0,
	MaxLogLevel: len(LogLevelString),
}

var LogLevelString = [...]string{
	"info",
	"trace",
	"debug",
	"warning",
	"error",
	"fatal",
}

func (m *model) LevelIntToString() {
	m.LogLevel++
	id := m.LogLevel % m.MaxLogLevel
	if m.LogLevel >= m.MaxLogLevel {
		m.LogLevel = 0
	}

	switch id {
	case 0: //nolint:goconst
		logrus.SetLevel(logrus.InfoLevel)
	case 1: //nolint:goconst
		logrus.SetLevel(logrus.TraceLevel)
	case 2: //nolint:goconst
		logrus.SetLevel(logrus.DebugLevel)
	case 3:
		logrus.SetLevel(logrus.WarnLevel)
	case 4: //nolint:goconst
		logrus.SetLevel(logrus.ErrorLevel)
	case 5: //nolint:goconst
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	fmt.Println("logLevel", LogLevelString[id])
}

// func Keyboard() (err error) {
// 	if err = keyboard.Open(); err != nil {
// 		return
// 	}
//
// 	defer func() {
// 		_ = keyboard.Close()
// 	}()
//
// FOR0:
// 	for {
// 		char, key, err := keyboard.GetKey()
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		switch key {
// 		case keyboard.KeyCtrlC:
// 			break FOR0
// 		}
//
// 		switch char {
// 		case 'q':
// 			break FOR0
// 		case 'c':
// 			break FOR0
// 		case 'l':
// 			mod.LevelIntToString()
// 		}
//
// 		if key == keyboard.KeyEsc {
// 			break FOR0
// 		}
// 	}
// 	//os.Exit(0)
// 	return
// }

func Add(cancel context.CancelFunc, translator intf.Translator) {
	readRU := false

	ctrl_c_func := func() {
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
			// translator.OnlyOriginal()
			translator.OnlyTranslate()
			translator.Go(processedString)
		}
	}

	// xclip
	go func() {
		messageChannel := make(chan string)
		go clipboardMonitor(messageChannel)

		for {
			select {
			case msg := <-messageChannel:
				fmt.Println("Received message:", msg)
				ctrl_c_func()
			case <-time.After(time.Second):
				// Do something else or just wait
			}
		}
	}()

	translator.OnlyOriginal()

	fmt.Println("--- Please press ctrl + q to stop hook ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl"}, func(e hook.Event) {
		fmt.Println("ctrl-q")
		translator.OnlyOriginalRu()
		translator.Go("завершение программы")
		time.Sleep(1 * time.Second)
		cancel()
	})

	hook.Register(hook.KeyDown, []string{"p", "ctrl", "alt"}, func(e hook.Event) {
		fmt.Println("ctrl-alt-p")

		if !translator.CheckPause() {
			translator.OnlyOriginalRu()
			translator.Go("пауза")
		} else {
			translator.OnlyOriginalRu()
			translator.Go("пауза снята")
		}
		translator.SetPause()
	})

	//hook.Register(hook.KeyDown, []string{"t", "alt"}, func(e hook.Event) {
	//fmt.Println("alt-t")
	////voice.InvertTranslate()

	////if voice.TanslateOrNot() {
	////translator.OnlyOriginalRu("без перевода")
	////} else {
	////translator.OnlyOriginalRu("переводить текст")
	////}
	//})

	//hook.Register(hook.KeyDown, []string{"-", "alt"}, func(e hook.Event) {
	//fmt.Println("-", "alt")
	////out, speed, err := voice.SpeedSub()
	////if err != nil {
	////fmt.Println(err)
	////return
	////}

	////logrus.WithFields(logrus.Fields{
	////"out": out,
	////}).Info("speed-")

	////str := fmt.Sprintf("%.1f", speed)
	////translator.OnlyOriginalRu(str)
	//})

	//hook.Register(hook.KeyDown, []string{"+", "alt"}, func(e hook.Event) {
	//fmt.Println("+", "alt")
	////out, speed, err := voice.SpeedAdd()
	////if err != nil {
	////fmt.Println(err)
	////return
	////}

	////logrus.WithFields(logrus.Fields{
	////"out": out,
	////}).Info("speed+")

	////str := fmt.Sprintf("%.1f", speed)
	////translator.OnlyOriginalRu(str)
	//})

	hook.Register(hook.KeyDown, []string{"f", "alt"}, func(e hook.Event) {
		if translator.CheckPause() {
			return
		}
		readRU = !readRU

		if readRU {
			translator.OnlyOriginalRu()
			translator.Go("режим чтения на русском")
		} else {
			translator.OnlyOriginalRu()
			translator.Go("включить переводчик")
		}
		// time.Sleep(time.Millisecond * 50)
		// text, err := clipboard.ReadAll()
		//
		// if err != nil {
		// 	logrus.WithFields(logrus.Fields{
		// 		"err": err,
		// 	}).Warn("clipboard")
		//
		// 	translator.OnlyOriginalRu()
		// 	translator.Go("не скопировалось")
		// 	return
		// }

		// processedString, err := RegexWork(text)

		// if err != nil {
		// 	logrus.WithFields(logrus.Fields{
		// 		"err": err,
		// 	}).Warn("regexp")
		// 	return
		// }
		// translator.TranslateAndOriginal()
		// translator.Go(processedString)
	})

	//hook.Register(hook.KeyDown, []string{"meta", "t"}, func(e hook.Event) {
	//fmt.Println("key down", e.String(), e.Keychar)
	//})

	fmt.Println("--- Please press t---")
	hook.Register(hook.KeyDown, []string{"z", "ctrl"}, func(e hook.Event) {
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
		// translator.OnlyTranslate()
		translator.Go(processedString)
	})

	//fmt.Println("--- Please press c---")
	//hook.Register(hook.KeyDown, []string{"z", "ctrl"}, func(e hook.Event) {
	//if translator.CheckPause() {
	//return
	//}

	// time.Sleep(time.Millisecond * 50)
	// text, err := clipboard.ReadAll()

	//if err != nil {
	//logrus.WithFields(logrus.Fields{
	//"err": err,
	//}).Warn("clipboard")

	//translator.OnlyOriginalRu()
	//translator.Go("не скопировалось")
	//return
	//}

	// processedString, err := RegexWork(text)

	//if err != nil {
	//logrus.WithFields(logrus.Fields{
	//"err": err,
	//}).Warn("regexp")
	//return
	//}
	//translator.OnlyOriginal()
	//translator.Go(processedString)
	//})

	fmt.Println("--- Please press c---")
	hook.Register(hook.KeyDown, []string{"c", "ctrl"}, func(e hook.Event) {
		ctrl_c_func()
	})

	hook.Register(hook.KeyDown, []string{"r", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("r", "ctrl", "shift")
	})

	s := hook.Start()
	<-hook.Process(s)
}

func Low() {
	EvChan := hook.Start()
	defer hook.End()

	for ev := range EvChan {
		fmt.Println("hook: ", ev)
	}
}

func Event() {
}

var (
	reg0             = regexp.MustCompile(`[^a-zA-Z\p{Han}0-9 .,\r\n]+`)
	reg2             = regexp.MustCompile(`([\p{L}])\.([\p{L}])`)
	reg3             = regexp.MustCompile(`([[:lower:]])([[:upper:]])`)
	reg4             = regexp.MustCompile(`(\b(\p{L}+)\b)`)
	singleSpaceRegex = regexp.MustCompile(`\s+`)

	re = regexp.MustCompile(`([\wλ]+)\s+∈\s+([\w ]+)`)

	mathSymbols = map[string]string{
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
	tt = reg01.ReplaceAllString(tt, " ")
	tt = singleSpacePattern.ReplaceAllString(tt, " ")
	tt = strings.ReplaceAll(tt, " .", ".")
	tt = strings.ReplaceAll(tt, " ,", ",")

	return strings.TrimSpace(tt), nil
}
