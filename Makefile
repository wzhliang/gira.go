.PHONY: build install

build:
	go build


install: build
	cp gira ~/bin/gira

