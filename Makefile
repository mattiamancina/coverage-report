# Makefile for parse_cobertura

BINARY_NAME = parse_cobertura
INPUT = input.xml
OUTPUT = coverage.md
CHANGELIST = changeList.txt

.PHONY: all build run clean

all: build run

build:
	go build -o $(BINARY_NAME) parse_cobertura.go

run: build
	./$(BINARY_NAME) -input $(INPUT) -output $(OUTPUT) -changeList $(CHANGELIST)

clean:
	rm -f $(BINARY_NAME) $(OUTPUT)
