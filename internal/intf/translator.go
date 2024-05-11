package intf

import "context"

type Translator interface {
	Run()
	OnlyTranslate()
	OnlyOriginal()
	OnlyOriginalRu()
	TranslateAndOriginal()
	Go(text string)
	CheckPause() bool
	SetPause()
	Speak(ctx context.Context, text, command string) string
}
