.PHONY: default
default:
	@echo 'This does nothing by default. Maybe you want to make install?'

.PHONY: install
install:
	go build -o /usr/local/bin/mind-meld .
