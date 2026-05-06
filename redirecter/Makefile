BINARY_NAME=redirecter
MAIN_PATH=./cmd/proxy/main.go

ifeq ($(OS),Windows_NT)
    RM = if exist bin rmdir /s /q bin
    EXT = .exe
    # Na Windowsie w Makefile zmienne dla Go ustawiamy tak:
    SET_LINUX = set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&&
    SET_WIN   = set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=0&&
else
    RM = rm -rf bin
    EXT =
    SET_LINUX = GOOS=linux GOARCH=amd64 CGO_ENABLED=0
    SET_WIN   = GOOS=windows GOARCH=amd64 CGO_ENABLED=0
endif

.PHONY: all build-linux build-windows clean

all: build-linux build-windows

build-linux:
	@echo --- Budowanie na Linux ---
	$(SET_LINUX) go build -o bin/$(BINARY_NAME)_linux $(MAIN_PATH)

build-windows:
	@echo --- Budowanie na Windows ---
	$(SET_WIN) go build -o bin/$(BINARY_NAME)$(EXT) $(MAIN_PATH)

clean:
	@echo --- Czyszczenie ---
	$(RM)
