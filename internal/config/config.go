package config

import (
	evdev "github.com/back2nix/golang-evdev"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LangSpeechConfig holds speed and half settings for a language.
type LangSpeechConfig struct {
	Speed int `mapstructure:"Speed"`
	Half  int `mapstructure:"Half"`
}

// SpeechConfig содержит настройки, связанные с синтезом речи для разных языков.
type SpeechConfig struct {
	DefaultOutput string           `mapstructure:"DefaultOutput"`
	En            LangSpeechConfig `mapstructure:"En"`
	Ru            LangSpeechConfig `mapstructure:"Ru"`
}

// Config определяет все параметры конфигурации для приложения.
type Config struct {
	Speech SpeechConfig      `mapstructure:"Speech"`
	Input  InputConfig       `mapstructure:"Input"`
	Sounds map[string]string `mapstructure:"Sounds"`
}

// InputConfig содержит настройки, связанные с устройствами ввода.
type InputConfig struct {
	Device  string        `mapstructure:"Device"`
	Hotkeys HotkeysConfig `mapstructure:"Hotkeys"`
	Listen  string        `mapstructure:"Listen"`
}

// HotkeysConfig определяет действия и соответствующие им комбинации клавиш в виде строк.
type HotkeysConfig struct {
	Translate        string `mapstructure:"Translate"`
	TranslateOral    string `mapstructure:"TranslateOral"`
	ToggleReadMode   string `mapstructure:"ToggleReadMode"`
	TogglePause      string `mapstructure:"TogglePause"`
	StopSound        string `mapstructure:"StopSound"`
	ToggleCopyBuffer string `mapstructure:"ToggleCopyBuffer"`
}

// KeyToCode сопоставляет строковые имена клавиш с их кодами evdev.
var KeyToCode = map[string]uint16{
	"C": evdev.KEY_C,
	"Z": evdev.KEY_Z,
	"F": evdev.KEY_F,
	"P": evdev.KEY_P,
	"V": evdev.KEY_V,
	// Добавьте другие клавиши по необходимости
}

// Init загружает конфигурацию из файла или устанавливает значения по умолчанию.
func Init() (*Config, error) {
	// Значения по умолчанию
	viper.SetDefault("Speech.DefaultOutput", "Translate")
	viper.SetDefault("Speech.En.Speed", 3)
	viper.SetDefault("Speech.En.Half", 2)
	viper.SetDefault("Speech.Ru.Speed", 7)
	viper.SetDefault("Speech.Ru.Half", 2)
	viper.SetDefault("Input.Device", "/dev/input/event1")
	viper.SetDefault("Input.Listen", ":3111")

	// ИСПОЛЬЗУЕМ КЛЮЧИ В НИЖНЕМ РЕГИСТРЕ
	viper.SetDefault("Sounds", map[string]string{
		"start":      "sound/interface-soft-click-131438.mp3",
		"processing": "sound/computer-processing.mp3",
		"click":      "sound/slide-click-92152.mp3",
	})

	// Горячие клавиши
	viper.SetDefault("Input.Hotkeys.Translate", "Ctrl+C")
	viper.SetDefault("Input.Hotkeys.TranslateOral", "Ctrl+Z")
	viper.SetDefault("Input.Hotkeys.ToggleReadMode", "Alt+F")
	viper.SetDefault("Input.Hotkeys.TogglePause", "Ctrl+Alt+P")
	viper.SetDefault("Input.Hotkeys.StopSound", "Alt+C")
	viper.SetDefault("Input.Hotkeys.ToggleCopyBuffer", "Alt+V")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Viper автоматически обрабатывает регистр ключей в YAML, приводя их к нижнему.
	// Для доступа к ним нужно использовать нижний регистр.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Warn("Config file not found. Using default values.")
		} else {
			return nil, err
		}
	} else {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal config")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"start_sound":      cfg.Sounds["start"],
		"processing_sound": cfg.Sounds["processing"],
		"click_sound":      cfg.Sounds["click"],
	}).Debug("Final loaded sound configuration")

	return &cfg, nil
}
