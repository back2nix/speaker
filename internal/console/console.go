package console

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/back2nix/speaker/internal/intf"
	hook "github.com/robotn/gohook"
	"github.com/sirupsen/logrus"
)

type model struct {
	logLevel    int
	maxLogLevel int
}

var logLevelString = [...]string{
	"info",
	"trace",
	"debug",
	"warning",
	"error",
	"fatal",
}

func (m *model) levelToIntString() {
	m.logLevel++
	id := m.logLevel % len(logLevelString)
	if m.logLevel >= len(logLevelString) {
		m.logLevel = 0
	}

	switch id {
	case 0:
		logrus.SetLevel(logrus.InfoLevel)
	case 1:
		logrus.SetLevel(logrus.TraceLevel)
	case 2:
		logrus.SetLevel(logrus.DebugLevel)
	case 3:
		logrus.SetLevel(logrus.WarnLevel)
	case 4:
		logrus.SetLevel(logrus.ErrorLevel)
	case 5:
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	fmt.Println("logLevel", logLevelString[id])
}

func add(cancel context.CancelFunc, translator intf.Translator) {
	readRU := false

	ctrlCFunc := func() {
		if translator.CheckPause() {
			return
		}

		time.Sleep(50 * time.Millisecond)
		text, err := clipboard.ReadAll()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("clipboard")

			translator.OnlyOriginalRu()
			translator.Go("не скопировалось")
			return
		}

		processedString, err := regexWork(text)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("regexp")
			return
		}
		if readRU {
			processedString, _ := regexWorkRu(text)
			translator.OnlyOriginalRu()
			translator.Go(processedString)
		} else {
			translator.OnlyTranslate()
			translator.Go(processedString)
		}
	}

	go func() {
		messageChannel := make(chan string)
		go clipboardMonitor(messageChannel)

		for {
			select {
			case msg := <-messageChannel:
				fmt.Println("Received message:", msg)
				ctrlCFunc()
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
	})

	hook.Register(hook.KeyDown, []string{"z", "ctrl"}, func(e hook.Event) {
		if translator.CheckPause() {
			return
		}

		time.Sleep(50 * time.Millisecond)
		text, err := clipboard.ReadAll()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("clipboard")

			translator.OnlyOriginalRu()
			translator.Go("не скопировалось")
			return
		}

		processedString, err := regexWork(text)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Warn("regexp")
			return
		}
		translator.OnlyOriginal()
		translator.Go(processedString)
	})

	hook.Register(hook.KeyDown, []string{"c", "ctrl"}, func(e hook.Event) {
		ctrlCFunc()
	})

	hook.Start()
	<-hook.Process()
}

func low() {
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

	return regexp.MustCompile(regexPattern)
}

func MathRegex(input string) (out string, err error) {
	out = regexPattern.ReplaceAllStringFunc(input, func(matched string) string {
		return mathSymbols[matched]
	})

	return out, nil
}

func regexWork(tt string) (out string, err error) {
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

func regexWorkRu(tt string) (out string, err error) {
	tt = reg01.ReplaceAllString(tt, " ")
	tt = singleSpacePattern.ReplaceAllString(tt, " ")
	tt = strings.ReplaceAll(tt, " .", ".")
	tt = strings.ReplaceAll(tt, " ,", ",")

	return strings.TrimSpace(tt), nil
}
