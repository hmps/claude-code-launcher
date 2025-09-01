.PHONY: build clean run

build:
	go build -o cc-launcher

clean:
	rm -f cc-launcher

run:
	go run main.go