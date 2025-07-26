package translateshell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/back2nix/speaker/internal/config"
	"github.com/sirupsen/logrus"
)

type Store struct {
	ctx           context.Context
	ctxSpeak      context.Context
	cancelSpeak   context.CancelFunc
	chText        chan string
	typeOperation string
	pause         bool
	original      string
	translate     string
	lastText      string
	speechConfig  *config.SpeechConfig
}

func New(ctx context.Context, speechCfg *config.SpeechConfig) (store *Store) {
	store = &Store{
		ctx:          ctx,
		chText:       make(chan string),
		speechConfig: speechCfg,
	}
	return
}

func (s *Store) Run() {
	s.ctxSpeak, s.cancelSpeak = context.WithCancel(s.ctx)
	for {
		var text string
		select {
		case <-s.ctx.Done():
			return
		case text = <-s.chText:
		}

		s.cancelSpeak()
		time.Sleep(20 * time.Millisecond)
		s.ctxSpeak, s.cancelSpeak = context.WithCancel(s.ctx)

		go func() {
			if text != s.lastText || s.lastText == "" {
				text = strings.ToLower(text)
				s.translate = speak(s.ctxSpeak, text, `trans -b -t ru "%s"`)
				s.lastText = text
				s.original = text
			}

			var err error
			switch s.typeOperation {
			case operationOnlyTranslate:
				speed := s.speechConfig.Ru.Speed
				half := s.speechConfig.Ru.Half
				err = replay(s.ctxSpeak, "ru", s.translate, speed, half)
				if err != nil {
					logrus.WithError(err).Error("Replay 'OnlyTranslate' failed")
				}
			case operationOnlyOriginalRu:
				speed := s.speechConfig.Ru.Speed
				half := s.speechConfig.Ru.Half
				err = replay(s.ctxSpeak, "ru", s.original, speed, half)
				if err != nil {
					logrus.WithError(err).Error("Replay 'OnlyOriginalRu' failed")
				}
			case operationOnlyOriginal:
				speed := s.speechConfig.En.Speed
				half := s.speechConfig.En.Half
				err = replay(s.ctxSpeak, "en", s.original, speed, half)
				if err != nil {
					logrus.WithError(err).Error("Replay 'OnlyOriginal' failed")
				}
			case operationTranslateAndOriginal:
				speedRu := s.speechConfig.Ru.Speed
				halfRu := s.speechConfig.Ru.Half
				err = replay(s.ctxSpeak, "ru", s.translate, speedRu, halfRu)
				if err != nil {
					logrus.WithError(err).Error("Replay 'TranslateAndOriginal' (RU) failed")
				}
				speedEn := s.speechConfig.En.Speed
				halfEn := s.speechConfig.En.Half
				err = replay(s.ctxSpeak, "en", s.original, speedEn, halfEn)
				if err != nil {
					logrus.WithError(err).Error("Replay 'TranslateAndOriginal' (EN) failed")
				}
			}
		}()
	}
}

const (
	operationOnlyTranslate        string = "OnlyTranslate"
	operationOnlyOriginal         string = "OnlyOriginal"
	operationOnlyOriginalRu       string = "OnlyOriginalRu"
	operationTranslateAndOriginal string = "TranslateAndOriginal"
)

func (s *Store) OnlyTranslate() {
	s.typeOperation = operationOnlyTranslate
}

func (s *Store) OnlyOriginal() {
	s.typeOperation = operationOnlyOriginal
}

func (s *Store) OnlyOriginalRu() {
	s.typeOperation = operationOnlyOriginalRu
}

func (s *Store) TranslateAndOriginal() {
	s.typeOperation = operationTranslateAndOriginal
}

func (s *Store) Go(text string) {
	s.chText <- text
}

func (s *Store) Speak(ctx context.Context, text, command string) string {
	return speak(ctx, text, command)
}

func (s *Store) CheckPause() bool {
	return s.pause
}

func (s *Store) SetPause() {
	s.pause = !s.pause
	if s.pause {
		s.cancelSpeak()
	}
}

func speak(ctx context.Context, text, command string) string {
	txtCmd := fmt.Sprintf(command, text)
	cmd := exec.CommandContext(ctx, "sh", "-c", txtCmd)
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	return string(out)
}

var c2 *exec.Cmd

func replay(ctx context.Context, lang, text string, speed, half int) (err error) {
	if c2 != nil && c2.Process != nil {
		_ = c2.Process.Kill()
	}

	if text == "" {
		return
	}

	strCommand := fmt.Sprintf(`gtts-cli -l %s "%s"`, lang, text)
	fmt.Println(text)
	c1 := exec.CommandContext(ctx, "bash", "-c", strCommand)
	stdout1, err := c1.StdoutPipe()
	if err != nil {
		logrus.WithError(err).Error("gtts-cli stdout pipe failed")
		return
	}
	if err = c1.Start(); err != nil {
		logrus.WithError(err).Error("gtts-cli start failed")
		return
	}

	strCommand2 := fmt.Sprintf(`mpg123 -d %d -h %d --pitch 0 -`, speed, half)
	c2 = exec.CommandContext(ctx, "bash", "-c", strCommand2)
	c2.Stdin = stdout1
	if err = c2.Start(); err != nil {
		logrus.WithError(err).Error("mpg123 start failed")
		return
	}
	if err = c1.Wait(); err != nil {
		logrus.WithError(err).Warn("gtts-cli wait finished with error")
	}
	if err = c2.Wait(); err != nil {
		logrus.WithError(err).Warn("mpg123 wait finished with error")
	}

	return nil
}

func Stop() {
	if c2 != nil && c2.Process != nil {
		_ = c2.Process.Kill()
	}
}

func Play(file string) {
	log := logrus.WithField("file", file)
	log.Debug("Play function called")

	if c2 != nil && c2.Process != nil {
		log.Debug("Stopping previous playback")
		_ = c2.Process.Kill()
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.WithError(err).Error("Sound file does not exist at the given path.")
		return
	}

	// Используем -q (quiet), чтобы подавить информационный вывод mpg123
	strCommand := fmt.Sprintf(`mpg123 -q "%s"`, file)
	log.WithField("command", strCommand).Info("Executing sound playback command")

	// Мы не используем CommandContext, т.к. звук должен играть в фоне без привязки к контексту операции
	cmd := exec.Command("sh", "-c", strCommand)

	go func() {
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.WithError(err).WithField("output", string(output)).Error("Sound playback command failed")
			return
		}
		if len(output) > 0 {
			log.WithField("output", string(output)).Debug("Sound playback command finished with output")
		} else {
			log.Debug("Sound playback command finished successfully")
		}
	}()
}
