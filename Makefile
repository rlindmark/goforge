PROGRAM=goforge

all: test build

test:
	go test .

build:
	go build -o $(PROGRAM)

clean:
	rm $(PROGRAM)
