# Makefile

all: build

build:
	go build -o terraform-provider-atlas.exe	

.PHONY: all build

