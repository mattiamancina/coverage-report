# Cobertura Coverage Parser

Cobertura Coverage Parser is a GitHub Action that parses a Cobertura XML coverage report and generates a Markdown summary of overall and per-file coverage.

## Features

- Parses Cobertura XML reports.
- Generates an line coverage summary for the changed files
- Generates an overall line summary.
- Lists individual file coverage in a collapsible section.
- Outputs a Markdown file for use in GitHub workflows.

## Inputs

| Name         | Description                          | Default      |
|--------------|--------------------------------------|--------------|
| `input`      | Path to Cobertura XML report file    | `input.xml`  |
| `changeList` | Path to a changelist file            |   |
| `output`     | Path to the generated Markdown report | `coverage.md`|

## Usage

Add this action to your GitHub workflow to parse a Cobertura XML report and produce a Markdown coverage report:

```yaml
name: Coverage Report

on: [push, pull_request]

jobs:
  coverage:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Generate Cobertura report
        run: |
          # Generate coverage.xml in Cobertura format
          # e.g., using jacoco, pytest-cov, or other tools

      - name: Parse Cobertura coverage
        uses: mattiamancina/coverage-report@v1.0.0
        with:
          input: coverage.xml
          changeList: changelist.txt
          output: coverage.md

      - name: Display Coverage Report
        run: cat coverage.md
```

## Implementation

This action is implemented in Go. The included Dockerfile compiles the Go program and packages it into a minimal Alpine-based container for use in GitHub Actions.

## Local Usage

Build and run the action locally with Docker:

```bash
docker build -t cobertura-parser .
docker run --rm -v "$(pwd)":/github/workspace cobertura-parser \
  -input coverage.xml \
  -output coverage.md
```

Or build the Go binary directly:

```bash
go build -o parse_cobertura parse_cobertura.go
./parse_cobertura -input coverage.xml -output coverage.md
```