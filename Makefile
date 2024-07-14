# Построить образ прежде чем собрать под windows
#docker_prebuild_image:
	#docker build -t back2nix/golang_bakend_msys2 -f docker/cross/Dockerfile .

#windows:
	##GOOS=windows GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -v -o build/speaker.exe cmd/voice/main.go
	#./scripts/cross.sh

run: build/speaker
	./build/speaker

.PHONY: build/speaker
build/speaker:
	go build -v -o build/speaker cmd/speaker/main.go
	#CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -a -installsuffix cgo -v -o build/android_speaker cmd/android/main.go

# not maintain
# .PHONY: keyloger_run
# keyloger_run:
# 	go build -o build/keylogger cmd/keylogger/keylogger.go && sudo ./build/keylogger

# not tested
# сборка deb пакета с использование flake
# deb:
# 	nix bundle --bundler bundlers#toDEB .

# not fully works
# appimage:
# 	nix bundle --bundler github:ralismark/nix-appimage .

# Просто сборка
nix/build:
	nix build .

lorri:
	systemctl --user daemon-reload
	systemctl --user status lorri.service
	lorri daemon
