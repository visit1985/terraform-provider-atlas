# Makefile

all: clean build

build:
	go build -o terraform-provider-atlas.exe
	terraform init
	terraform plan

clean:
	rm -f terraform-provider-atlas.exe

.PHONY: all build clean

