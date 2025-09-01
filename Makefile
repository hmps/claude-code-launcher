.PHONY: build clean run

build:
	go build -o cc-launcher

clean:
	rm -f cc-launcher

install:
	go install

run:
	go run main.go
