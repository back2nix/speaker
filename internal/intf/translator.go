package intf

type Translator interface {
	Run()
	OnlyTranslate()
	OnlyOriginal()
	OnlyOriginalRu()
	TranslateAndOriginal()
	Go(text string)
	CheckPause() bool
	SetPause()
}
