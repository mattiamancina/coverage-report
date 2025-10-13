package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

// CoverageXML represents the root cobertura XML element
type CoverageXML struct {
	LineRate   float64      `xml:"line-rate,attr"`
	BranchRate float64      `xml:"branch-rate,attr"`
	Packages   []PackageXML `xml:"packages>package"`
}

// PackageXML represents a cobertura package element
type PackageXML struct {
	Classes []ClassXML `xml:"classes>class"`
}

// ClassXML represents a cobertura class element
type ClassXML struct {
	Filename   string  `xml:"filename,attr"`
	LineRate   float64 `xml:"line-rate,attr"`
	BranchRate float64 `xml:"branch-rate,attr"`
}

// FileCoverage holds coverage data for a single file
type FileCoverage struct {
	Filename   string
	LineRate   float64
	BranchRate float64
}

func generateBadge(coverage float64) string {
	if coverage >= 75 {
		return "green"
	}

	if coverage > 70 && coverage < 75 {
		return "yellow"
	}
	return "red"
}

func main() {
	inputPath := flag.String("input", "input.xml", "Path to Cobertura XML report file")
	changeListPath := flag.String("changeList", "i", "Path to change list file")
	outputPath := flag.String("output", "coverage.md", "Path to output Markdown file")
	flag.Parse()

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading XML file '%s': %v\n", *inputPath, err)
		os.Exit(1)
	}

	changeListFile, err := os.Open(*changeListPath)

	var changedFiles []string
	scanner := bufio.NewScanner(changeListFile)
	for scanner.Scan() {
		changedFiles = append(changedFiles, strings.TrimPrefix(scanner.Text(), "api/app/"))
	}

	var cov CoverageXML
	if err := xml.Unmarshal(data, &cov); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing XML: %v\n", err)
		os.Exit(1)
	}

	var changedFilesCoverage []FileCoverage

	var avgCoverage float64
	// Collect file coverage
	var files []FileCoverage
	for _, pkg := range cov.Packages {
		for _, cls := range pkg.Classes {
			files = append(files, FileCoverage{
				Filename:   cls.Filename,
				LineRate:   cls.LineRate,
				BranchRate: cls.BranchRate,
			})

			if len(changedFiles) < 1 {
				continue
			}

			if slices.Contains(changedFiles, cls.Filename) {
				changedFilesCoverage = append(changedFilesCoverage, FileCoverage{
					Filename:   cls.Filename,
					LineRate:   cls.LineRate,
					BranchRate: cls.BranchRate,
				})
				avgCoverage += cls.LineRate
			}
		}
	}

	// Sort by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].Filename < files[j].Filename
	})

	// Build Markdown report
	var b strings.Builder
	b.WriteString("# Coverage Report\n\n")

	b.WriteString("\n")
	b.WriteString("\n")

	if len(changedFilesCoverage) > 0 {

		b.WriteString(fmt.Sprintf("**Changed Files Line Coverage:** %.2f%%  \n", avgCoverage/float64(len(changedFilesCoverage))*100))
		b.WriteString("\n")
		b.WriteString("\n")
		b.WriteString("![badge](https://img.shields.io/badge/Coverage-" + fmt.Sprintf("%f", avgCoverage/float64(len(changedFilesCoverage))*100) + "%25-" + generateBadge(avgCoverage/float64(len(changedFilesCoverage))*100) + ")")
		b.WriteString("\n")
		b.WriteString("\n")
		b.WriteString("| File | Line Coverage | \n")
		b.WriteString("| ---- | ------------- | \n")
		for _, f := range changedFilesCoverage {
			b.WriteString(fmt.Sprintf("| %s | %.2f%% |\n", f.Filename, f.LineRate*100))
		}
	}
	b.WriteString("\n")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("**Overall Line Coverage:** %.2f%%  \n", cov.LineRate*100))
	b.WriteString("\n")
	b.WriteString("![badge](https://img.shields.io/badge/Coverage-" + fmt.Sprintf("%f", cov.LineRate*100) + "%25-" + generateBadge(cov.LineRate*100) + ")")
	b.WriteString("\n")
	b.WriteString("<details>")
	b.WriteString("\n")
	b.WriteString("<summary>File Coverage</summary>\n")
	b.WriteString("\n")
	b.WriteString("| File | Line Coverage | \n")
	b.WriteString("| ---- | ------------- | \n")
	for _, f := range files {
		b.WriteString(fmt.Sprintf("| %s | %.2f%% |\n", f.Filename, f.LineRate*100))
	}

	b.WriteString("</details>")

	// Write to output file
	if err := os.WriteFile(*outputPath, []byte(b.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file '%s': %v\n", *outputPath, err)
		os.Exit(1)
	}
}
