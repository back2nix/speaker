package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MarinX/keylogger"
	"github.com/sirupsen/logrus"

	key "github.com/back2nix/speaker/internal/keylogger"
)

func main() {
	channel := make(chan string)

	go KeyMonitor(channel, key.KeyCombos)

	for key := range channel {
		fmt.Printf("Received key event: %s\n", key)
		SendMessage(key)
	}
}

func SendMessage(input string) {
	url := fmt.Sprintf("http://localhost:3111/echo/%s", input)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Response:", resp.Status)
	} else {
		fmt.Println("Error:", resp.Status)
	}
}

// KeyMonitor функция мониторинга клавиш
func KeyMonitor(channel chan string, keyCombos [][]string) {
	// find keyboard device, does not require a root permission
	keyboard := keylogger.FindKeyboardDevice()

	// check if we found a path to keyboard
	if len(keyboard) <= 0 {
		logrus.Error("No keyboard found...you will need to provide manual input path")
		return
	}

	logrus.Println("Found a keyboard at", keyboard)
	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer k.Close()

	events := k.Read()

	var (
		lastKeys      []string
		lastShortsCut string
	)

	for e := range events {
		if e.Type == keylogger.EvKey {
			keyString := e.KeyString()

			if e.KeyPress() {
				lastKeys = append(lastKeys, keyString)
				if len(lastKeys) > 3 {
					lastKeys = lastKeys[1:]
				}

				stringMap := make(map[string]bool)

				// Добавляем элементы среза в map
				for _, item := range lastKeys {
					stringMap[item] = true
				}

				for _, keyCombo := range keyCombos {
					if len(lastKeys) == 3 && len(keyCombo) == 3 &&
						lastKeys[0] == keyCombo[0] &&
						lastKeys[1] == keyCombo[1] &&
						lastKeys[2] == keyCombo[2] {

						shortsCut := strings.Join(keyCombo, "+")
						if shortsCut == "L_ALT+C" && lastShortsCut == shortsCut {
							channel <- strings.Join(keyCombo, "+") + "x2"
							lastShortsCut = ""
						} else {
							channel <- shortsCut
							lastShortsCut = shortsCut
						}
					}
					if len(lastKeys) == 3 && len(keyCombo) == 2 &&
						lastKeys[1] == keyCombo[0] &&
						lastKeys[2] == keyCombo[1] {
						shortsCut := strings.Join(keyCombo, "+")
						if shortsCut == "L_ALT+C" && lastShortsCut == shortsCut {
							channel <- strings.Join(keyCombo, "+") + "x2"
							lastShortsCut = ""
						} else {
							channel <- shortsCut
							lastShortsCut = shortsCut
						}
					}

					if len(lastKeys) == 3 &&
						lastKeys[0] == "L_ALT" &&
						lastKeys[1] == "C" {
						if lastShortsCut == "L_ALT+C" {
							channel <- "L_ALT+Cx2"
							lastShortsCut = ""
						}
					}
				}
			}
		}
	}
}

func containsElement(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
