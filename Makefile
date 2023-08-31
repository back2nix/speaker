# Построить образ прежде чем собрать под windows
#docker_prebuild_image:
	#docker build -t back2nix/golang_bakend_msys2 -f docker/cross/Dockerfile .

#windows:
	##GOOS=windows GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -v -o build/speaker.exe cmd/voice/main.go
	#./scripts/cross.sh

# dependencis
# https://github.com/pndurette/gTTS
# https://github.com/soimort/translate-shell

run: build/speaker
	./build/speaker

install_depend:
	go install github.com/playwright-community/playwright-go/cmd/playwright
	playwright install --with-deps
	sudo -H pip3 install gTTS
	#go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps

.PHONY: build/speaker
build/speaker:
	go build -v -o build/speaker cmd/speaker/main.go
	#CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -a -installsuffix cgo -v -o build/android_speaker cmd/android/main.go

.PHONY: keyloger_run
keyloger_run:
	go build -o build/keylogger cmd/keylogger/keylogger.go && sudo ./build/keylogger


# сборка deb пакета с использование flake
deb:
	nix bundle --bundler bundlers#toDEB .

appimage:
	nix bundle --bundler github:ralismark/nix-appimage .

# Просто сборка
nix/build:
	nix-build
