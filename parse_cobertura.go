package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
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

func main() {
	inputPath := flag.String("input", "input.xml", "Path to Cobertura XML report file")
	outputPath := flag.String("output", "coverage.md", "Path to output Markdown file")
	flag.Parse()

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading XML file '%s': %v\n", *inputPath, err)
		os.Exit(1)
	}

	var cov CoverageXML
	if err := xml.Unmarshal(data, &cov); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing XML: %v\n", err)
		os.Exit(1)
	}

	// Collect file coverage
	var files []FileCoverage
	for _, pkg := range cov.Packages {
		for _, cls := range pkg.Classes {
			files = append(files, FileCoverage{
				Filename:   cls.Filename,
				LineRate:   cls.LineRate,
				BranchRate: cls.BranchRate,
			})
		}
	}

	// Sort by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].Filename < files[j].Filename
	})

	// Build Markdown report
	var b strings.Builder
	b.WriteString("# Coverage Report\n\n")
	b.WriteString(fmt.Sprintf("**Overall Line Coverage:** %.2f%%  \n", cov.LineRate*100))
	b.WriteString(fmt.Sprintf("**Overall Branch Coverage:** %.2f%%  \n", cov.BranchRate*100))
	b.WriteString("\n")
	b.WriteString("<details>")
	b.WriteString("<summary>File Coverage</summary>\n")
	b.WriteString("\n")
	b.WriteString("| File | Line Coverage | Branch Coverage |\n")
	b.WriteString("| ---- | ------------- | --------------- |\n")
	for _, f := range files {
		b.WriteString(fmt.Sprintf("| %s | %.2f%% | %.2f%% |\n", f.Filename, f.LineRate*100, f.BranchRate*100))
	}

	b.WriteString("</details>")

	// Write to output file
	if err := os.WriteFile(*outputPath, []byte(b.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file '%s': %v\n", *outputPath, err)
		os.Exit(1)
	}
}
